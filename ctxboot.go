package ctxboot

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

var (
	defaultContext *ComponentContext
	once           sync.Once
)

// Boot returns the default component context
func Boot() *ComponentContext {
	once.Do(func() {
		defaultContext = NewComponentContext()
	})
	return defaultContext
}

// ComponentContext manages components and their dependencies
type ComponentContext struct {
	components map[reflect.Type]interface{}
	mu         sync.RWMutex
}

// NewComponentContext creates a new component context
func NewComponentContext() *ComponentContext {
	return &ComponentContext{
		components: make(map[reflect.Type]interface{}),
	}
}

// GetComponent retrieves a component by its type
func (cc *ComponentContext) GetComponent(typ reflect.Type) (interface{}, error) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	// First try exact match
	if component, ok := cc.components[typ]; ok {
		return component, nil
	}

	// If the requested type is an interface, look for implementations
	if typ.Kind() == reflect.Interface {
		var candidates []reflect.Type
		for t := range cc.components {
			// Check if the component type implements the interface
			if t.Implements(typ) {
				candidates = append(candidates, t)
			}
		}

		// If no candidates found, return error
		if len(candidates) == 0 {
			return nil, fmt.Errorf("no component found that implements interface: %v", typ)
		}

		// If multiple candidates found, panic
		if len(candidates) > 1 {
			panic(fmt.Sprintf("multiple components implement interface %v: %v", typ, candidates))
		}

		// Return the single candidate
		return cc.components[candidates[0]], nil
	}

	return nil, fmt.Errorf("component not found: %v", typ)
}

// SetComponent stores a component instance
func (cc *ComponentContext) SetComponent(typ reflect.Type, instance interface{}) error {
	if instance == nil {
		return errors.New("cannot store nil component")
	}

	if !reflect.TypeOf(instance).AssignableTo(typ) {
		return fmt.Errorf("instance type %v is not assignable to %v", reflect.TypeOf(instance), typ)
	}

	// Check if component already exists
	cc.mu.RLock()
	if _, exists := cc.components[typ]; exists {
		cc.mu.RUnlock()
		return nil // Skip if component already exists
	}
	cc.mu.RUnlock()

	// Store the component
	cc.mu.Lock()
	cc.components[typ] = instance
	cc.mu.Unlock()

	return nil
}

// InitializeComponents injects dependencies into all registered components
func (cc *ComponentContext) InitializeComponents() error {
	// Create a copy of components to avoid concurrent modification
	components := make(map[reflect.Type]interface{})
	cc.mu.RLock()
	for typ, comp := range cc.components {
		components[typ] = comp
	}
	cc.mu.RUnlock()

	// Track initialized components
	initialized := make(map[reflect.Type]bool)

	// Initialize components until all are done or we can't make progress
	for len(initialized) < len(components) {
		progress := false

		for typ, instance := range components {
			if initialized[typ] {
				continue
			}

			// Check if all dependencies are initialized
			val := reflect.ValueOf(instance)
			if val.Kind() != reflect.Ptr {
				return fmt.Errorf("component must be a pointer: %v", typ)
			}

			elem := val.Elem()
			if elem.Kind() != reflect.Struct {
				return fmt.Errorf("component must be a pointer to a struct: %v", typ)
			}

			allDepsInitialized := true
			for i := 0; i < elem.Type().NumField(); i++ {
				field := elem.Type().Field(i)
				if tag := field.Tag.Get("ctxboot"); tag == "inject" {
					fieldType := field.Type
					if fieldType.Kind() != reflect.Ptr {
						fieldType = reflect.PtrTo(fieldType)
					}

					if _, exists := components[fieldType]; exists && !initialized[fieldType] {
						allDepsInitialized = false
						break
					}
				}
			}

			if allDepsInitialized {
				if err := cc.injectDependencies(instance); err != nil {
					return fmt.Errorf("failed to initialize component %v: %w", typ, err)
				}
				initialized[typ] = true
				progress = true
			}
		}

		if !progress {
			// Find uninitialized components for error message
			var uninitialized []string
			for typ := range components {
				if !initialized[typ] {
					uninitialized = append(uninitialized, typ.String())
				}
			}
			return fmt.Errorf("circular dependency detected among: %v", uninitialized)
		}
	}

	return nil
}

// injectDependencies injects dependencies into a component
func (cc *ComponentContext) injectDependencies(target interface{}) error {
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}

	elem := val.Elem()
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to a struct")
	}

	typ := elem.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if tag := field.Tag.Get("ctxboot"); tag == "inject" {
			// Get the type for the field
			fieldType := field.Type
			isPtrField := fieldType.Kind() == reflect.Ptr

			// For interface fields, use the interface type directly
			if fieldType.Kind() == reflect.Interface {
				component, err := cc.GetComponent(fieldType)
				if err != nil {
					return fmt.Errorf("failed to inject field %s: %w", field.Name, err)
				}

				fieldVal := elem.Field(i)
				if !fieldVal.CanSet() {
					// Handle unexported field
					fieldVal = reflect.NewAt(field.Type, unsafe.Pointer(fieldVal.UnsafeAddr())).Elem()
				}

				// Set the value
				fieldVal.Set(reflect.ValueOf(component))
				continue
			}

			// For non-interface fields
			lookupType := fieldType
			if !isPtrField {
				lookupType = reflect.PtrTo(fieldType)
			}

			component, err := cc.GetComponent(lookupType)
			if err != nil {
				return fmt.Errorf("failed to inject field %s: %w", field.Name, err)
			}

			fieldVal := elem.Field(i)
			if !fieldVal.CanSet() {
				// Handle unexported field
				fieldVal = reflect.NewAt(field.Type, unsafe.Pointer(fieldVal.UnsafeAddr())).Elem()
			}

			// Convert component to the correct type
			compVal := reflect.ValueOf(component)
			if !isPtrField {
				// If field is not a pointer, dereference the component
				compVal = compVal.Elem()
			}

			// Set the value
			fieldVal.Set(compVal)
		}
	}
	return nil
}

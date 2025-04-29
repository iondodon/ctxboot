package ctxboot

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

// CtxbootComponentContext manages components and their dependencies
type CtxbootComponentContext struct {
	components map[reflect.Type]interface{}
	mu         sync.RWMutex
}

// NewCtxbootComponentContext creates a new component context
func NewCtxbootComponentContext() *CtxbootComponentContext {
	return &CtxbootComponentContext{
		components: make(map[reflect.Type]interface{}),
	}
}

// GetComponent retrieves a component by its type
func (c *CtxbootComponentContext) GetComponent(typ reflect.Type) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// First try exact match
	if component, ok := c.components[typ]; ok {
		return component, nil
	}

	// If the requested type is an interface, look for implementations
	if typ.Kind() == reflect.Interface {
		var candidates []reflect.Type
		for t := range c.components {
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
		return c.components[candidates[0]], nil
	}

	return nil, fmt.Errorf("component not found: %v", typ)
}

// SetComponent stores a component instance
func (c *CtxbootComponentContext) SetComponent(typ reflect.Type, instance interface{}) error {
	if instance == nil {
		return errors.New("cannot store nil component")
	}

	// Get the actual type of the instance
	instanceType := reflect.TypeOf(instance)

	// If typ is a pointer type but instance is not, create a pointer to instance
	if typ.Kind() == reflect.Ptr && instanceType.Kind() != reflect.Ptr {
		// Create a new pointer to the instance
		ptr := reflect.New(instanceType)
		ptr.Elem().Set(reflect.ValueOf(instance))
		instance = ptr.Interface()
		instanceType = reflect.TypeOf(instance)
	}

	if !instanceType.AssignableTo(typ) {
		return fmt.Errorf("instance type %v is not assignable to %v", instanceType, typ)
	}

	// Check if component already exists
	c.mu.RLock()
	if _, exists := c.components[typ]; exists {
		c.mu.RUnlock()
		return nil // Skip if component already exists
	}
	c.mu.RUnlock()

	// Store the component
	c.mu.Lock()
	c.components[typ] = instance
	c.mu.Unlock()

	return nil
}

// InitializeComponents injects dependencies into all registered components
func (c *CtxbootComponentContext) InitializeComponents() error {
	// Create a copy of components to avoid concurrent modification
	components := make(map[reflect.Type]interface{})
	c.mu.RLock()
	for typ, comp := range c.components {
		components[typ] = comp
	}
	c.mu.RUnlock()

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
				if err := c.injectDependencies(instance); err != nil {
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
func (c *CtxbootComponentContext) injectDependencies(target interface{}) error {
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
				component, err := c.GetComponent(fieldType)
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

			component, err := c.GetComponent(lookupType)
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

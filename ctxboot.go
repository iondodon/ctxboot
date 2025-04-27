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

	if component, ok := cc.components[typ]; ok {
		return component, nil
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

	// Initialize each component
	for typ, instance := range components {
		if err := cc.injectDependencies(instance); err != nil {
			return fmt.Errorf("failed to initialize component %v: %w", typ, err)
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

			// If field is not a pointer, we need to get the pointer type to look up the component
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

			// If field is not a pointer, dereference the component
			if !isPtrField {
				component = reflect.ValueOf(component).Elem().Interface()
			}

			fieldVal.Set(reflect.ValueOf(component))
		}
	}
	return nil
}

// MustRegister registers a component and panics if registration fails
func MustRegister(component interface{}) {
	typ := reflect.TypeOf(component)
	if typ.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("component must be a pointer, got %v", typ))
	}

	if err := Boot().SetComponent(typ, component); err != nil {
		panic(fmt.Sprintf("failed to register component: %v", err))
	}
}

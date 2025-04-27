package ctxboot

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
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

	// Store the component first to handle circular dependencies
	cc.mu.Lock()
	cc.components[typ] = instance
	cc.mu.Unlock()

	// Initialize the component
	if err := cc.injectDependencies(instance); err != nil {
		cc.mu.Lock()
		delete(cc.components, typ) // Rollback on error
		cc.mu.Unlock()
		return fmt.Errorf("failed to inject dependencies: %v", err)
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
			component, err := cc.GetComponent(field.Type)
			if err != nil {
				return fmt.Errorf("failed to inject field %s: %w", field.Name, err)
			}

			fieldVal := elem.Field(i)
			if !fieldVal.CanSet() {
				return fmt.Errorf("cannot set field %s", field.Name)
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

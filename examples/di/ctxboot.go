// Code generated by ctxboot; DO NOT EDIT.

package main

import (
	"github.com/iondodon/ctxboot"
	"reflect"
	"log"
	"fmt"
	
	database "github.com/iondodon/ctxboot/examples/di/database"
	
	repository "github.com/iondodon/ctxboot/examples/di/repository"
	
)

// ComponentContext embeds CtxbootComponentContext and adds getter methods
type ComponentContext struct {
	*ctxboot.CtxbootComponentContext
}

// RegisterComponent registers a component instance and automatically deduces its type
func (c *ComponentContext) RegisterComponent(instance interface{}) error {
	if instance == nil {
		return fmt.Errorf("cannot register nil component")
	}
	return c.SetComponent(reflect.TypeOf(instance), instance)
}

// registerScanedComponenets registers all components
func (c *ComponentContext) registerScanedComponenets() error {
	// Register components in dependency order
	
	// Register database.DatabaseImpl
	if err := c.SetComponent(reflect.TypeOf((*database.DatabaseImpl)(nil)), &database.DatabaseImpl{}); err != nil {
		log.Fatalf("Failed to register component %s: %v", "database.DatabaseImpl", err)
	}
	
	// Register UserService
	if err := c.SetComponent(reflect.TypeOf((*UserService)(nil)), &UserService{}); err != nil {
		log.Fatalf("Failed to register component %s: %v", "UserService", err)
	}
	
	// Register repository.UserRepository
	if err := c.SetComponent(reflect.TypeOf((*repository.UserRepository)(nil)), &repository.UserRepository{}); err != nil {
		log.Fatalf("Failed to register component %s: %v", "repository.UserRepository", err)
	}
	
	
	return nil
}

// NewComponentContext creates a new component context instance and registers all scanned components
func NewComponentContext() *ComponentContext {
	ctx := &ComponentContext{ctxboot.NewCtxbootComponentContext()}
	if err := ctx.registerScanedComponenets(); err != nil {
		log.Fatalf("Failed to register scanned components: %v", err)
	}
	return ctx
}

// Component getter methods

// GetDatabaseImpl returns the DatabaseImpl component
func (c *ComponentContext) GetDatabaseImpl() (*database.DatabaseImpl, error) {
	component, err := c.GetComponent(reflect.TypeOf((*database.DatabaseImpl)(nil)))
	if err != nil {
		return nil, err
	}
	return component.(*database.DatabaseImpl), nil
}

// GetUserService returns the UserService component
func (c *ComponentContext) GetUserService() (*UserService, error) {
	component, err := c.GetComponent(reflect.TypeOf((*UserService)(nil)))
	if err != nil {
		return nil, err
	}
	return component.(*UserService), nil
}

// GetUserRepository returns the UserRepository component
func (c *ComponentContext) GetUserRepository() (*repository.UserRepository, error) {
	component, err := c.GetComponent(reflect.TypeOf((*repository.UserRepository)(nil)))
	if err != nil {
		return nil, err
	}
	return component.(*repository.UserRepository), nil
}


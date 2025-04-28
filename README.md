# CtxBoot

A lightweight dependency injection framework for Go.

## Usage

### 1. Define Components

Add the `ctxboot:component` annotation to your structs:

```go
//ctxboot:component
type Database struct {
    // Database implementation
}

//ctxboot:component
type UserRepository struct {
    DB *Database `ctxboot:"inject"`
}

//ctxboot:component
type UserService struct {
    Repo *UserRepository `ctxboot:"inject"`
}
```

### 2. Run the Code Generator

```bash
go run cmd/ctxboot/main.go ./path/to/your/package
```

This will generate a `ctxboot.go` file with component registration code.

### 3. Use in Your Application

```go
package main

import (
    "log"
    "github.com/iondodon/ctxboot"
)

func main() {
    // Register custom components before loading context
    cc := ctxboot.Boot()
    if err := cc.SetComponent(reflect.TypeOf((*CustomComponent)(nil)), &CustomComponent{}); err != nil {
        log.Fatalf("Failed to register custom component: %v", err)
    }

    // Load and initialize all components
    if err := LoadContext(cc); err != nil {
        log.Fatalf("Failed to load context: %v", err)
    }

    // Get a component instance
    service, err := cc.GetComponent(reflect.TypeOf((*UserService)(nil)))
    if err != nil {
        log.Fatalf("Failed to get service: %v", err)
    }

    // Use the component
    userService := service.(*UserService)
    // ... use userService
}
```

Note: If you need to register custom components, do it before calling `LoadContext`. Components registered after `LoadContext` will not have their dependencies injected.

## Features

- Automatic dependency injection
- Component lifecycle management
- Support for both pointer and non-pointer components
- Thread-safe component context
- Dependency order initialization
- Circular dependency detection

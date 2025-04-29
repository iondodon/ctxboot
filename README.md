# Ctxboot

Ctxboot is a lightweight dependency injection framework for Go that helps manage component lifecycle and dependencies.

## Features

- Automatic dependency injection
- Component lifecycle management
- Type-safe component access
- Support for both interface and concrete type dependencies
- Circular dependency detection
- Thread-safe operations
- Support for unexported fields

## Components

The framework consists of two main components:

1. **CtxbootComponentContext** (Library)

   - Core dependency injection container
   - Manages component registration and lifecycle
   - Handles dependency resolution and injection

2. **ComponentContext** (Generated)
   - Application-specific context
   - Embeds CtxbootComponentContext
   - Provides type-safe getter methods for components
   - Adds application-specific functionality

## Usage

### 1. Define Components

Mark your components with the `ctxboot:component` annotation:

```go
// ctxboot:component
type MyComponent struct {
    Dependency *OtherComponent `ctxboot:"inject"`
}
```

### 2. Generate Code

Run the code generator:

```bash
go run cmd/ctxboot/main.go <package-dir>
```

This will generate a `ctxboot.go` file with:

- Component registration code
- Type-safe getter methods
- Context initialization code

### 3. Use in Your Application

```go
// Create a new context
ctx := NewComponentContext()

// Register components
if err := ctx.RegisterScanedComponenets(); err != nil {
    log.Fatal(err)
}

// Initialize components and inject dependencies
if err := ctx.InjectComponents(); err != nil {
    log.Fatal(err)
}

// Get components
myComp, err := ctx.GetMyComponent()
if err != nil {
    log.Fatal(err)
}
```

## Example

```go
// ctxboot:component
type Database struct {
    Config *Config `ctxboot:"inject"`
}

// ctxboot:component
type Config struct {
    // configuration fields
}

// ctxboot:component
type Service struct {
    DB *Database `ctxboot:"inject"`
}

func main() {
    ctx := NewComponentContext()

    if err := ctx.RegisterScanedComponenets(); err != nil {
        log.Fatal(err)
    }

    if err := ctx.InjectComponents(); err != nil {
        log.Fatal(err)
    }

    service, err := ctx.GetService()
    if err != nil {
        log.Fatal(err)
    }

    // Use the service...
}
```

## Best Practices

1. **Component Naming**

   - Use clear, descriptive names for components
   - Components must be exported (start with capital letter)

2. **Dependency Management**

   - Keep dependency graphs shallow
   - Avoid circular dependencies
   - Use interfaces for better testability

3. **Error Handling**

   - Always check errors from context operations
   - Handle initialization failures gracefully

4. **Thread Safety**
   - The framework is thread-safe
   - Components should be thread-safe if accessed concurrently

## License

MIT License

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
- Component overriding (later registrations replace earlier ones)
- Automatic component registration

## Components

The framework consists of two main components:

1. **CtxbootComponentContext** (Library)

   - Core dependency injection container
   - Manages component registration and lifecycle
   - Handles dependency resolution and injection
   - Supports component overriding

2. **ComponentContext** (Generated)
   - Application-specific context
   - Embeds CtxbootComponentContext
   - Provides type-safe getter methods for components
   - Adds application-specific functionality
   - Automatically registers scanned components on creation

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
// Create a new context (automatically registers scanned components)
ctx := NewComponentContext()

// You can override components by registering them again
customComponent := &MyComponent{/* ... */}
if err := ctx.RegisterComponent(customComponent); err != nil {
    log.Fatal(err)
}

// Initialize components and inject dependencies
if err := ctx.InitializeComponents(); err != nil {
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
    // Create context (automatically registers all scanned components)
    ctx := NewComponentContext()

    // Override the default Database with a mock for testing
    mockDB := &Database{/* mock implementation */}
    if err := ctx.RegisterComponent(mockDB); err != nil {
        log.Fatal(err)
    }

    if err := ctx.InitializeComponents(); err != nil {
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
   - Be aware that later component registrations override earlier ones

3. **Error Handling**

   - Always check errors from context operations
   - Handle initialization failures gracefully

4. **Thread Safety**

   - The framework is thread-safe
   - Components should be thread-safe if accessed concurrently

5. **Component Overriding**
   - Use component overriding for testing (replacing real components with mocks)
   - Be careful with component overriding in production code
   - Document when components are meant to be overridden
   - Consider using interfaces to make component overriding more predictable

## License

MIT License

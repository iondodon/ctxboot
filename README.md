# Ctxboot

Ctxboot is a dependency injection library for Go that helps manage component dependencies.

See [https://github.com/iondodon/ctxboot/tree/main/examples](https://github.com/iondodon/ctxboot/tree/main/examples)

## Usage

### 1. Define Components

Mark your components with the `ctxboot:component` annotation:

```go
//ctxboot:component
type MyComponent struct {
    Dependency *OtherComponent `ctxboot:"inject"`
}
```

### 2. Generate Code

Run the code generator:

```bash
go install github.com/iondodon/ctxboot/cmd/ctxboot@latest
```

in your project

```bash
ctxboot .
```

or

```bash
ctxboot <relative directory path>
```

This will generate a `ctxboot.go` file with:

- Component registration code
- Type-safe getter methods
- Context initialization code

### 3. Use in Your Application

```go
// Create a new context (automatically registers scanned components)
cc := NewComponentContext()

// REgister custom configured components
// with RegisterComponent you can override existing components
customComponent := &MyComponent{/* custom cinfigureation */}
if err := cc.RegisterComponent(customComponent); err != nil {
    log.Fatal(err)
}

// Initialize components and inject dependencies
if err := cc.InitializeComponents(); err != nil {
    log.Fatal(err)
}

// Get components
myComp, err := cc.GetMyComponent()
if err != nil {
    log.Fatal(err)
}
```

## Example

```go
//ctxboot:component
type Database struct {
    Config *Config `ctxboot:"inject"`
}

type Config struct {
    // configuration fields
}

//ctxboot:component
type Service struct {
    DB *Database `ctxboot:"inject"`
}

func main() {
    // Create context (automatically registers all scanned components)
    cc := NewComponentContext()

    config := &Config{/* configuration */}
    if err := cc.RegisterComponent(config); err != nil {
        log.Fatal(err)
    }

    if err := cc.InitializeComponents(); err != nil {
        log.Fatal(err)
    }

    service, err := cc.GetService()
    if err != nil {
        log.Fatal(err)
    }

    // Use the service...
}
```

## License

Apache License Version 2.0, January 2004

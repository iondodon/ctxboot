# CtxBoot

A lightweight dependency injection framework for Go.

## Usage

1. Add the `ctxboot:component` annotation to your struct:

```go
//ctxboot:component
type MyComponent struct {
    // ...
}
```

2. Add the `ctxboot:"inject"` tag to fields that should be injected:

```go
type MyComponent struct {
    Dependency *OtherComponent `ctxboot:"inject"`
}
```

3. Run the code generator:

```bash
go run cmd/ctxboot/main.go ./path/to/your/package
```

4. Load the context in your main function:

```go
func main() {
    if err := LoadContext(ctxboot.Boot()); err != nil {
        log.Fatalf("Failed to load context: %v", err)
    }

    // Your application code...
}
```

## Features

- Automatic dependency injection
- Component lifecycle management
- Support for both pointer and non-pointer components
- Thread-safe component context

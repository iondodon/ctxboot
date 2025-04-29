# CtxBoot

A lightweight dependency injection framework for Go.

## Installation

```bash
go install github.com/iondodon/ctxboot/cmd/ctxboot@latest
```

## Example

Here's a comprehensive example that demonstrates all key features of CtxBoot:

```go
// Define custom components (not scanned)
type LoggerConfig struct {
    Prefix string
    Flags  int
}

type DatabaseConfig struct {
    Host     string
    Port     int
    Username string
    Password string
}

// Define scanned components
//ctxboot:component
type Database struct {
    Config *DatabaseConfig `ctxboot:"inject"`
}

//ctxboot:component
type UserRepository struct {
    DB     *Database      `ctxboot:"inject"`
    Logger *LoggerConfig  `ctxboot:"inject"`
}

//ctxboot:component
type UserService struct {
    Repo   *UserRepository `ctxboot:"inject"`
    Logger *LoggerConfig   `ctxboot:"inject"`
}

func main() {
    // Create a new context
    cc := NewContext()

    // Register custom components first
    loggerConfig := &LoggerConfig{
        Prefix: "APP: ",
        Flags:  log.Ldate | log.Ltime,
    }
    if err := cc.RegisterComponent(loggerConfig); err != nil {
        log.Fatal(err)
    }

    dbConfig := &DatabaseConfig{
        Host:     "localhost",
        Port:     5432,
        Username: "postgres",
        Password: "secret",
    }
    if err := cc.RegisterComponent(dbConfig); err != nil {
        log.Fatal(err)
    }

    // Then register scanned components
    if err := cc.RegisterScanedComponenets(); err != nil {
        log.Fatal(err)
    }

    // Initialize components and their dependencies
    if err := cc.InjectComponents(); err != nil {
        log.Fatal(err)
    }

    // Get components using generated getter methods (only for scanned components)
    userService, err := cc.GetUserService()
    if err != nil {
        log.Fatal(err)
    }

    // Get components by type (works for both scanned and non-scanned components)
    logger, err := cc.GetComponent(reflect.TypeOf((*LoggerConfig)(nil)))
    if err != nil {
        log.Fatal(err)
    }
    loggerConfig := logger.(*LoggerConfig)

    // Get scanned components by type (alternative to generated getters)
    db, err := cc.GetComponent(reflect.TypeOf((*Database)(nil)))
    if err != nil {
        log.Fatal(err)
    }
    database := db.(*Database)

    // Use the components
    user, err := userService.GetUser("123")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("User: %+v\n", user)
}
```

This example demonstrates:

1. **Component Definition**

   - Custom components without `//ctxboot:component` annotation
   - Scanned components with `//ctxboot:component` annotation
   - Dependency injection using `ctxboot:"inject"` tag

2. **Component Registration**

   - Manual registration of custom components using `RegisterComponent`
   - Automatic registration of scanned components using `RegisterScanedComponenets`
   - Proper order of registration (custom components first)

3. **Dependency Injection**

   - Automatic injection of both custom and scanned components
   - Support for pointer and non-pointer types
   - Dependency resolution in correct order

4. **Component Access**
   - Type-safe access using generated getter methods (only for scanned components)
   - Generic access using `GetComponent` (works for all components)
   - Error handling for missing components

Note: Only components marked with `//ctxboot:component` will have generated getter methods. For components without this annotation, you must use `GetComponent` to retrieve them.

## License

Apache License Version 2.0, January 2004

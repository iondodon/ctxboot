package main

import (
	"fmt"
	"reflect"
)

// Define an interface
type Greeter interface {
	Greet() string
}

// Implement the interface
//
//ctxboot:component
type EnglishGreeter struct{}

func (g *EnglishGreeter) Greet() string {
	return "Hello!"
}

func main() {
	// Create a new context
	cc := NewComponentContext()

	// Register components
	if err := cc.RegisterScanedComponenets(); err != nil {
		panic(err)
	}

	// Initialize components and their dependencies
	if err := cc.InjectComponents(); err != nil {
		panic(err)
	}

	// Example 1: Get component by interface type
	greeter, err := cc.GetComponent(reflect.TypeOf((*Greeter)(nil)).Elem())
	if err != nil {
		panic(err)
	}
	g := greeter.(Greeter)
	fmt.Println("Example 1 - Get by interface:")
	fmt.Println(g.Greet())

	// Example 2: Get component using generated getter method
	englishGreeter, err := cc.GetEnglishGreeter()
	if err != nil {
		panic(err)
	}
	fmt.Println("Example 2 - Get by generated getter:")
	fmt.Println(englishGreeter.Greet())
}

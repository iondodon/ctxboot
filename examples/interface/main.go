package main

import (
	"fmt"
	"reflect"

	"github.com/iondodon/ctxboot"
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
	cc := ctxboot.Boot()
	err := LoadContext(cc)
	if err != nil {
		panic(err)
	}

	// Get component by interface type
	greeter, err := cc.GetComponent(reflect.TypeOf((*Greeter)(nil)).Elem())
	if err != nil {
		panic(err)
	}

	// Use the component
	g := greeter.(Greeter)
	fmt.Println(g.Greet())
}

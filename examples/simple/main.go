package main

import (
	"fmt"
	"reflect"

	"github.com/iondodon/ctxboot"
)

// UserService handles user-related operations
//
//ctxboot:component
type UserService struct {
	// Add fields as needed
}

func (s *UserService) GetUser(id string) string {
	return fmt.Sprintf("User %s", id)
}

func main() {
	// Get component
	service, err := ctxboot.Boot().GetComponent(reflect.TypeOf(&UserService{}))
	if err != nil {
		panic(err)
	}

	// Use component
	userService := service.(*UserService)
	fmt.Println(userService.GetUser("123"))
}

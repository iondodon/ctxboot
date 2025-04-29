package main

import (
	"fmt"
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
	cc, err := LoadContext()
	if err != nil {
		panic(err)
	}

	// Get component using the generated getter method
	userService, err := cc.GetUserService()
	if err != nil {
		panic(err)
	}

	// Use component
	fmt.Println(userService.GetUser("123"))
}

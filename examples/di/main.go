package main

import (
	"fmt"

	"github.com/iondodon/ctxboot/examples/di/repository"
)

// UserService provides user-related business logic
//
//ctxboot:component
type UserService struct {
	Repo *repository.UserRepository `ctxboot:"inject"`
}

func (s *UserService) GetUser(id string) string {
	return s.Repo.GetUser(id)
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
	fmt.Println("Example - Get by generated getter method:")
	fmt.Println(userService.GetUser("123"))
}

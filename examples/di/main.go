package main

import (
	"fmt"
	"reflect"

	"github.com/iondodon/ctxboot"
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
	// Get service component
	service, err := ctxboot.Boot().GetComponent(reflect.TypeOf(&UserService{}))
	if err != nil {
		panic(err)
	}

	// Use service
	userService := service.(*UserService)
	fmt.Println(userService.GetUser("123"))
}

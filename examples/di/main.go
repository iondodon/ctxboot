package main

import (
	"fmt"
	"reflect"

	"github.com/iondodon/ctxboot"
	"github.com/iondodon/ctxboot/examples/di/database"
)

// UserRepository handles user data access
//
//ctxboot:component
type UserRepository struct {
	DB *database.Database `ctxboot:"inject"`
}

func (r *UserRepository) GetUser(id string) string {
	if r.DB.ConnectionString == "" {
		r.DB.Connect()
	}
	return fmt.Sprintf("User %s from DB: %s", id, r.DB.ConnectionString)
}

// UserService provides user-related business logic
//
//ctxboot:component
type UserService struct {
	Repo *UserRepository `ctxboot:"inject"`
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

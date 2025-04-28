package main

import (
	"fmt"
	"reflect"

	"github.com/iondodon/ctxboot"
	"github.com/iondodon/ctxboot/examples/di/database"
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
	cc := ctxboot.Boot()
	err := LoadContext(cc)
	if err != nil {
		panic(err)
	}

	// Example 1: Get component by concrete type
	service, err := cc.GetComponent(reflect.TypeOf(&UserService{}))
	if err != nil {
		panic(err)
	}
	userService := service.(*UserService)
	fmt.Println("Example 1 - Get by concrete type:")
	fmt.Println(userService.GetUser("123"))

	// Example 2: Get component by interface - verbose way
	var dbInterface database.Database
	dbType := reflect.TypeOf(&dbInterface).Elem()
	db, err := cc.GetComponentByInterface(dbType)
	if err != nil {
		panic(err)
	}
	dbImpl := db.(database.Database)
	dbImpl.Connect()
	fmt.Println("\nExample 2 - Get by interface (verbose):")
	fmt.Printf("Database connection: %s\n", dbImpl.GetConnectionString())

	// Example 3: Get component by interface - concise way
	dbImpl2 := cc.MustGetComponentByInterface(reflect.TypeOf((*database.Database)(nil)).Elem()).(database.Database)
	dbImpl2.Connect()
	fmt.Println("\nExample 3 - Get by interface (concise):")
	fmt.Printf("Database connection: %s\n", dbImpl2.GetConnectionString())

	// Example 4: Get component by interface - with error handling
	dbImpl3, err := cc.GetComponentByInterface(reflect.TypeOf((*database.Database)(nil)).Elem())
	if err != nil {
		fmt.Printf("\nExample 4 - Error getting database: %v\n", err)
	} else {
		db := dbImpl3.(database.Database)
		db.Connect()
		fmt.Printf("\nExample 4 - Database connection: %s\n", db.GetConnectionString())
	}
}

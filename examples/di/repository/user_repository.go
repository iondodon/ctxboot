package repository

import (
	"fmt"

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

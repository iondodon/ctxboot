package repository

import (
	"fmt"

	"github.com/iondodon/ctxboot/examples/di/database"
)

// UserRepository handles user data access
//
//ctxboot:component
type UserRepository struct {
	db *database.Database `ctxboot:"inject"`
}

func (r *UserRepository) GetUser(id string) string {
	if r.db.ConnectionString == "" {
		r.db.Connect()
	}
	return fmt.Sprintf("User %s from DB: %s", id, r.db.ConnectionString)
}

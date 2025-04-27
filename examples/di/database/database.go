package database

// Database handles database operations
//
//ctxboot:component
type Database struct {
	ConnectionString string
}

func (db *Database) Connect() {
	db.ConnectionString = "connected"
}

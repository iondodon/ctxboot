package database

// Database interface defines database operations
type Database interface {
	Connect()
	GetConnectionString() string
}

// DatabaseImpl handles database operations
//
//ctxboot:component
type DatabaseImpl struct {
	ConnectionString string
}

func (db *DatabaseImpl) Connect() {
	db.ConnectionString = "connected"
}

func (db *DatabaseImpl) GetConnectionString() string {
	return db.ConnectionString
}

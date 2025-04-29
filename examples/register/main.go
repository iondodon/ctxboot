package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"
)

// LoggerConfig holds configuration for a logger
//
//ctxboot:component
type LoggerConfig struct {
	Prefix string
	Flags  int
}

// DatabaseConfig holds configuration for a database connection
//
//ctxboot:component
type DatabaseConfig struct {
	ConnectionString string
	MaxOpenConns     int
	MaxIdleConns     int
}

func main() {
	cc, err := LoadContext()
	if err != nil {
		log.Fatal(err)
	}

	// Register a LoggerConfig
	config := &LoggerConfig{
		Prefix: "APP: ",
		Flags:  log.Ldate | log.Ltime | log.Lshortfile,
	}
	if err := cc.RegisterComponent(config); err != nil {
		log.Fatal(err)
	}

	// Register a DatabaseConfig
	dbConfig := &DatabaseConfig{
		ConnectionString: "postgres://user:pass@localhost/db",
		MaxOpenConns:     10,
		MaxIdleConns:     5,
	}
	if err := cc.RegisterComponent(dbConfig); err != nil {
		log.Fatal(err)
	}

	// Register a time.Time (non-pointer type)
	startTime := time.Now()
	startTimePtr := &startTime
	if err := cc.RegisterComponent(startTimePtr); err != nil {
		log.Fatal(err)
	}

	// Register a *sql.DB (simulated)
	db := &sql.DB{} // In real code, this would be from sql.Open()
	if err := cc.RegisterComponent(db); err != nil {
		log.Fatal(err)
	}

	// Register a custom logger
	infoLog := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	if err := cc.RegisterComponent(infoLog); err != nil {
		log.Fatal(err)
	}

	// Initialize all components
	if err := cc.InitializeComponents(); err != nil {
		log.Fatal(err)
	}

	// Retrieve components
	loggerConfig, err := cc.GetLoggerConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Retrieved LoggerConfig: %+v\n", loggerConfig)

	dbConfig, err = cc.GetDatabaseConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Retrieved DatabaseConfig: %+v\n", dbConfig)

	// Get time component
	timeComponent, err := cc.GetComponent(reflect.TypeOf(&time.Time{}))
	if err != nil {
		log.Fatal(err)
	}
	startTime = *timeComponent.(*time.Time)
	fmt.Printf("Retrieved startTime: %v\n", startTime)

	dbInterface, err := cc.GetComponent(reflect.TypeOf(&sql.DB{}))
	if err != nil {
		log.Fatal(err)
	}
	db = dbInterface.(*sql.DB)
	fmt.Printf("Retrieved DB: %v\n", db)

	loggerInterface, err := cc.GetComponent(reflect.TypeOf(infoLog))
	if err != nil {
		log.Fatal(err)
	}
	logger := loggerInterface.(*log.Logger)
	fmt.Printf("Retrieved Logger: %v\n", logger)
}

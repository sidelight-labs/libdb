package repository

import (
	"fmt"
	"os"
)

const (
	DBHostEnv           = "DB_HOST"
	DBUserEnv           = "DB_USER"
	DBPasswordEnv       = "DB_PASSWORD"
	DBPortEnv           = "DB_PORT"
	RecordsFound        = "Found [%d] records to insert in block [%d]"
	DatabaseErrorPrefix = "Database Error:"
	MissingEnvError     = DatabaseErrorPrefix + " Missing the %s environment variable"
	PortError           = DatabaseErrorPrefix + " " + DBPortEnv + " must be an integer"
)

var env = []string{DBHostEnv, DBUserEnv, DBPasswordEnv, DBPortEnv}

type Store interface {
	Query(query string, args ...interface{}) (RowsScanner, error)
	Exec(statement string, args ...interface{}) error
}

type Scanner interface {
	Scan(dest ...interface{}) error
}

type RowsScanner interface {
	Next() bool
	Close() error
	Err() error
	Scanner
}

type config struct {
	host     string
	port     int
	user     string
	password string
	database string
}

func validateEnv() error {
	for _, value := range env {
		envVar := os.Getenv(value)
		if envVar == "" {
			return fmt.Errorf(MissingEnvError, value)
		}
	}
	return nil
}

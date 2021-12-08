package repository

import (
	"database/sql"
	"fmt"
	"github.com/sidelight-labs/libc/logger"
	"os"
	"strconv"
)

const (
	DBHostEnv       = "DB_HOST"
	DBUserEnv       = "DB_USER"
	DBPasswordEnv   = "DB_PASSWORD"
	DBPortEnv       = "DB_PORT"
	MySQLError      = "MySQL Error:"
	MissingEnvError = MySQLError + " Missing the %s environment variable"
	PortError       = MySQLError + " " + DBPortEnv + " must be an integer"
)

var env = []string{DBHostEnv, DBUserEnv, DBPasswordEnv, DBPortEnv}

type Store interface {
	Query(query string, args ...interface{}) (RowsScanner, error)
	Exec(statement string, args ...interface{}) (Summarizer, error)
}

type Scanner interface {
	Scan(dest ...interface{}) error
}

type RowsScanner interface {
	Columns() ([]string, error)
	Next() bool
	Close() error
	Err() error
	Scanner
}

type Summarizer interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

type config struct {
	host     string
	port     int
	user     string
	password string
	database string
}

type MySQL struct {
	config config
}

func NewMySQL(database string) (MySQL, error) {
	var result MySQL

	if err := validateEnv(); err != nil {
		return result, err
	}

	result.config.host = os.Getenv(DBHostEnv)
	result.config.user = os.Getenv(DBUserEnv)
	result.config.password = os.Getenv(DBPasswordEnv)
	result.config.database = database

	port, err := strconv.Atoi(os.Getenv(DBPortEnv))
	if err != nil {
		return result, fmt.Errorf(PortError)
	}
	result.config.port = port

	return result, nil
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

func (m MySQL) Query(query string, args ...interface{}) (RowsScanner, error) {
	db, err := sql.Open("mysql", m.connectionString())
	if err != nil {
		return nil, logger.Wrap(err, fmt.Sprintf("database error for query %s with args %v", query, args))
	}
	defer db.Close()

	return db.Query(query, args...)
}

func (m MySQL) Exec(statement string, args ...interface{}) (Summarizer, error) {
	db, err := sql.Open("mysql", m.connectionString())
	if err != nil {
		return nil, logger.Wrap(err, fmt.Sprintf("database error for statement %s with args %v", statement, args))
	}
	defer db.Close()

	return db.Exec(statement, args...)
}

func (m MySQL) connectionString() string {
	// <username>:<pw>@tcp(<HOST>:<port>)/<dbname>
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		m.config.user, m.config.password, m.config.host, m.config.port, m.config.database)
}

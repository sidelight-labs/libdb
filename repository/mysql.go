package repository

import (
	"database/sql"
	"fmt"
	"github.com/sidelight-labs/libc/logger"
	"os"
	"strconv"
)

const (
	MySQLError      = "MySQL Error:"
	MissingEnvError = MySQLError + " Missing the %s environment variable"
	PortError       = MySQLError + " " + DBPortEnv + " must be an integer"
)

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

func (m MySQL) Query(query string, args ...interface{}) (RowsScanner, error) {
	db, err := sql.Open("mysql", m.connectionString())
	if err != nil {
		return nil, logger.Wrap(err, fmt.Sprintf("database error for query %s with args %v", query, args))
	}
	defer db.Close()

	return db.Query(query, args...)
}

func (m MySQL) Exec(statement string, args ...interface{}) error {
	db, err := sql.Open("mysql", m.connectionString())
	if err != nil {
		return logger.Wrap(err, fmt.Sprintf("database error for statement %s with args %v", statement, args))
	}
	defer db.Close()

	_, err = db.Exec(statement, args...)

	return err
}

func (m MySQL) connectionString() string {
	// <username>:<pw>@tcp(<HOST>:<port>)/<dbname>
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		m.config.user, m.config.password, m.config.host, m.config.port, m.config.database)
}

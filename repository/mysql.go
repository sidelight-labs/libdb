package repository

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
)

type MySQL struct {
	config config
	db     *sql.DB
}

func NewMySQL(database string) (Store, error) {
	var result = &MySQL{}

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

	db, err := sql.Open("mysql", result.connectionString())
	if err != nil {
		return nil, err
	}

	result.db = db
	return result, nil
}

func (m *MySQL) Query(query string, args ...interface{}) (RowsScanner, error) {
	return m.db.Query(query, args...)
}

func (m *MySQL) Exec(statement string, args ...interface{}) error {
	_, err := m.db.Exec(statement, args...)
	return err
}

func (m *MySQL) Close() error {
	if m.db != nil {
		return m.db.Close()
	}

	return nil
}

func (m *MySQL) connectionString() string {
	// <username>:<pw>@tcp(<HOST>:<port>)/<dbname>
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		m.config.user, m.config.password, m.config.host, m.config.port, m.config.database)
}

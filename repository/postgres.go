package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4"
	"os"
	"strconv"
)

const (
	PostgresError     = "Postgres Error:"
	MissingEnvErrorPG = PostgresError + " Missing the %s environment variable"
	PortErrorPG       = PostgresError + " " + DBPortEnv + " must be an integer"
)

type Postgres struct {
	config config
}

func NewPostgres(database string) (Store, error) {
	var result Postgres

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

func (p Postgres) Query(query string, args ...interface{}) (RowsScanner, error) {
	// TODO: Not this
	panic("implement me")
}

func (p Postgres) Exec(statement string, args ...interface{}) error {
	conn, err := pgx.Connect(context.Background(), p.connectionString())
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), statement)
	conn.Query(context.Background(), statement)
	fmt.Printf(statement)
	return err
}

func (p Postgres) connectionString() string {
	// <username>:<pw>@tcp(<HOST>:<port>)/<dbname>
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		p.config.user, p.config.password, p.config.host, p.config.port, p.config.database)
}

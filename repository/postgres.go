package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4"
	"os"
	"strconv"
)

type Postgres struct {
	config config
	conn   *pgx.Conn
}

func NewPostgres(database string) (Store, error) {
	var result = &Postgres{}

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

	conn, err := pgx.Connect(context.Background(), result.connectionString())
	if err != nil {
		return nil, err
	}

	result.conn = conn
	return result, nil
}

func (p *Postgres) Query(query string, args ...interface{}) (RowsScanner, error) {
	// TODO: Not this
	panic("implement me")
}

func (p *Postgres) Exec(statement string, args ...interface{}) error {
	_, err := p.conn.Exec(context.Background(), statement)
	p.conn.Query(context.Background(), statement)
	fmt.Printf(statement)
	return err
}

func (p *Postgres) Close() error {
	if p.conn != nil {
		return p.conn.Close(context.Background())
	}

	return nil
}

func (p *Postgres) connectionString() string {
	// <username>:<pw>@tcp(<HOST>:<port>)/<dbname>
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		p.config.user, p.config.password, p.config.host, p.config.port, p.config.database)
}

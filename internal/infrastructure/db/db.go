package db

import (
	"database/sql"

	// Register PostgreSQL driver for database/sql.
	// This driver is required for sql.Open() to work with PostgreSQL.
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type IDB interface {
	GetterDBInterface
	Open() error
	Close() error
}

type GetterDBInterface interface {
	Db() *sql.DB
}

type PostgresSQLDB struct {
	db     *sql.DB
	params Params
}

func (p *PostgresSQLDB) Close() error {
	return p.Db().Close()
}

func (p *PostgresSQLDB) Db() *sql.DB {
	return p.db
}

func (p *PostgresSQLDB) Open() error {
	db, err := sql.Open(string(goose.DialectPostgres), p.params.DatabaseSourceName)

	if err != nil {
		return err
	}

	p.db = db

	return nil
}

func NewPostgresSQLDB(params Params) *PostgresSQLDB {
	return &PostgresSQLDB{params: params}
}

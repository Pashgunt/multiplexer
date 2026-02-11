package db

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

type DBInterface interface {
	GetterDBInterface
	Open() error
	Close() error
}

type GetterDBInterface interface {
	Db() *sql.DB
}

type PostgresSQLDB struct {
	db                 *sql.DB
	databaseSourceName string
}

func (p *PostgresSQLDB) Close() error {
	return p.Db().Close()
}

func (p *PostgresSQLDB) Db() *sql.DB {
	return p.db
}

func (p *PostgresSQLDB) Open() error {
	//postgres
	db, err := sql.Open(string(goose.DialectPostgres), p.databaseSourceName)

	if err != nil {
		return err
	}

	p.db = db

	return nil
}

func NewPostgresSQLDB(databaseSourceName string) *PostgresSQLDB {
	return &PostgresSQLDB{databaseSourceName: databaseSourceName}
}

package migrations

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

const migrationDir = "."

//go:embed *.sql
var EmbedMigrations embed.FS

type MigratorInterface interface {
	Setup() error
	Up() error
	Down() error
	Status() error
	Version() (int64, error)
	Reset() error
}

type Migrator struct {
	db *sql.DB
}

func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) Setup() error {
	goose.SetBaseFS(EmbedMigrations)

	return goose.SetDialect(string(goose.DialectPostgres))
}

func (m *Migrator) Up() error {
	return goose.Up(m.db, migrationDir)
}

func (m *Migrator) Down() error {
	return goose.Down(m.db, migrationDir)
}

func (m *Migrator) Status() error {
	return goose.Status(m.db, migrationDir)
}

func (m *Migrator) Version() (int64, error) {
	return goose.GetDBVersion(m.db)
}

func (m *Migrator) Reset() error {
	return goose.Reset(m.db, migrationDir)
}

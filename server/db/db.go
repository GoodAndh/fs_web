package db

import (
	"backend/config"
	"context"
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(env *config.Config) (*Database, error) {
	return databaseMySql(mysql.Config{
		User:                 env.DBUser,
		Passwd:               env.DBPassword,
		Net:                  "tcp",
		Addr:                 env.DBAddress,
		DBName:               validate(env),
		AllowNativePasswords: true,
		ParseTime:            true,
	})
}

func databaseMySql(cfg mysql.Config) (*Database, error) {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}

func validate(env *config.Config) string {

	if env.TestOrMain == "main" {
		return env.DBName
	} else {
		return env.DBTest
	}
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) DB() *sql.DB {
	return d.db
}

func (d *Database) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return d.db.ExecContext(ctx, query, args...)
}

func (d *Database) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return d.db.PrepareContext(ctx, query)
}

func (d *Database) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return d.db.QueryContext(ctx, query, args...)
}

func (d *Database) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return d.db.QueryRowContext(ctx, query, args...)
}

func (d *Database) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return d.db.BeginTx(ctx, opts)
}

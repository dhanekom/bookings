package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDBConn = 10
const maxIdleDbConn = 5
const maxDbLifetime = 5 & time.Minute

// ConnectSQL creates database pool for Postgres
func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	d.SetMaxOpenConns(maxOpenDBConn)
	d.SetConnMaxIdleTime(maxIdleDbConn)
	d.SetConnMaxLifetime(maxDbLifetime)

	dbConn.SQL = d

	err = testDB(d)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

// testDB tries to ping the database
func testDB(d *sql.DB) error {
	return d.Ping()
}

// NewDatabase creates a database for the application
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

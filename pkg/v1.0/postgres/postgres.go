package postgres

import (
	"database/sql"
	"fmt"

	"github.com/kenshaw/envcfg"
	_ "github.com/lib/pq"
)

type ResInsertInfo struct {
	LastInsertId int64
	RowsAffected int64
	ErrorMsg     string
}

const (
	POSTGRES      = `postgres`
	CUSTOM_HEALTH = `custom.health`
)

type PgConn string

// New initializes the postgres writers
func New(conf *envcfg.Envcfg) (*PgConn, error) {
	// Set postgres whitelist config
	pgc := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.GetKey("postgres.host"),
		conf.GetKey("postgres.port"),
		conf.GetKey("postgres.user"),
		conf.GetKey("postgres.password"),
		conf.GetKey("postgres.dbname"))

	db, err := sql.Open(POSTGRES, pgc)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(fmt.Sprintf(`SELECT * FROM %s LIMIT 1;`, CUSTOM_HEALTH))
	defer db.Close()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	pgconn := PgConn(pgc)

	return &pgconn, nil
}

package database

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewDBConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Configurações de pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Tenta conectar com retry
	for i := 0; i < 30; i++ {
		err = db.Ping()
		if err == nil {
			return db, nil
		}
		time.Sleep(1 * time.Second)
	}

	return nil, err
}

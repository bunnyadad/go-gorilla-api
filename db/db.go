package db

import (
	"database/sql"
	"fmt"
	"log"
)

type DB struct {
	Database *sql.DB
}

func (db *DB) Initialize(user, password, dbhost, dbname string) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, dbhost, dbname)

	var err error
	db.Database, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
}

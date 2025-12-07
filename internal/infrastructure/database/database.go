package database

import (
	"log"

	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB
var DRIVER string

func Get() *sqlx.DB {
	if DB == nil {
		log.Fatal("database connection is not initialized")
	}
	return DB
}

func GetDriver() string {
	return DRIVER
}

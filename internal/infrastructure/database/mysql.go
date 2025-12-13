package database

import (
	"fmt"
	"log"

	"github.com/fazriegi/fintrack-be/internal/infrastructure/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewMysql() {
	host := config.GetString("DB_HOST")
	username := config.GetString("DB_USERNAME")
	password := config.GetString("DB_PASSWORD")
	name := config.GetString("DB_NAME")
	port := config.GetInt("DB_PORT")
	DRIVER = config.GetString("DB_DRIVER")

	dbSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		username,
		password,
		host,
		port,
		name,
	)

	var err error
	DB, err = sqlx.Open(DRIVER, dbSource)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("failed to ping database:", err)
	}

	log.Println("connected to database successfully")
}

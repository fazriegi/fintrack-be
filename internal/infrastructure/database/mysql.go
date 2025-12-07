package database

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

func NewMysql(viper *viper.Viper) {
	host := viper.GetString("db.host")
	username := viper.GetString("db.username")
	password := viper.GetString("db.password")
	name := viper.GetString("db.name")
	port := viper.GetInt32("db.port")
	DRIVER = viper.GetString("db.driver")

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

package mysql

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
	
	const (
		DB_HOST = "49.247.134.77"
		DB_PORT = 3306
		DB_USERNAME = "sparker"
		DB_PWD = "tlchd50wh!"
		DB_DATABASE = "seoulbitz"
	)

	db, err := sql.Open("mysql", "sparker:tlchd50wh!@tcp(49.247.134.77:3306)/seoulbitz")
	if err != nil {
		panic(err)
	}

	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	
	return db
}
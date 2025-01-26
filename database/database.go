package database

import "database/sql"

var DB *sql.DB

func SetDB(db *sql.DB) {
	DB = db
}
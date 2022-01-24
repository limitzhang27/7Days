package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

const DB_DRIVER = "sqlite3"

func main() {
	db, _ := sql.Open(DB_DRIVER, "gee.db")
	defer func() { _ = db.Close() }()

	_, _ = db.Exec("DROP TABLE IF EXISTS User")
	_, _ = db.Exec("CREATE TABLE User(Name text, Age integer)")

	result, err := db.Exec("INSERT INTO User(Name, Age) VALUES ('Tom', 18), ('Jack', 25)")
	if err == nil {
		affected, _ := result.RowsAffected()
		log.Println(affected)
	}

	row := db.QueryRow("SELECT Name FROM User LIMIT 1")
	var name string
	if err := row.Scan(&name); err == nil {
		log.Println(name)
	}

}

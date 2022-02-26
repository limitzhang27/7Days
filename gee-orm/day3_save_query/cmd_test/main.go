package main

import (
	"fmt"
	"geeorm"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	engine, _ := geeorm.NewEngine("sqlite3", "gee.db")
	defer engine.Close()

	s := engine.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXITS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success, %d affeted\n", count)

	rows, _ := s.Raw("SELECT Name FROM User;").QueryRows()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err == nil {
			fmt.Printf("Name = %s\n", name)
		}
	}
}

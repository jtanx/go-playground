package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./data.sqlite")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = db.Exec(`
		--DROP TABLE IF EXISTS queries;
		CREATE TABLE IF NOT EXISTS queries (
			id integer primary key autoincrement,
			timestamp datetime default current_timestamp,
			parameters text default '' not null
		)
	`)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ret, err := db.Exec(`INSERT INTO queries (parameters) values (?)`, "lollolol; DROP TABLE queries;--")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	id, _ := ret.LastInsertId()
	fmt.Println("ID:", id)
	err = db.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("DONE")
}

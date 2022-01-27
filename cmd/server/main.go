package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dataSource := "root:@tcp(localhost:3306)/storage"
	var err error
	StorageDB, err := sql.Open("mysql", dataSource)
	if err != nil {
		log.Println(err)
	}

	if err = StorageDB.Ping(); err != nil {
		log.Println(err)
	}
	log.Println("Database Configured")
	StorageDB.Close()
}

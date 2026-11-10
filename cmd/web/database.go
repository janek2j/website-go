package main

import (
	"database/sql"
	"fmt"

	_ "github.com/glebarez/go-sqlite"
)

// var appCfg *config.AppConfig

func DBProcessFromAppConfig() {
	// db := config.AppConfig.Db
	// db := appCfg.Db
}

func DbProcess() {
	// Connect to the SQLite database
	db, err := sql.Open("sqlite", "./app.db")
	if err != nil {
		fmt.Println("error opening the database:", err)
		return
	}
	fmt.Println("Connected to the SQLite database successfully.")
	defer db.Close()

	// Get the version of SQLite
	var sqliteVersion string
	err = db.QueryRow("select sqlite_version()").Scan(&sqliteVersion)
	if err != nil {
		fmt.Println(err)
		return
	}
	db.Close()

	fmt.Println(sqliteVersion)

	db, err = sql.Open("sqlite", "app.db")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Db opened")
}

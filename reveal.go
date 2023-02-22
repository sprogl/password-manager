package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/sprogl/password-manager/database"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

func main() {
	dbFileName := "database.db"
	dbFileExists := fileExists(dbFileName)
	if !dbFileExists {
		log.Println("Creating database.db...")
		file, err := os.Create("database.db") // Create SQLite file
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		log.Println("database.db created")
	}

	db, _ := sql.Open("sqlite3", "./database.db") // Open the created SQLite File

	dbLogFile, err := os.OpenFile("db.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer dbLogFile.Close()

	dbHandler := database.CreateDBHandler(db, dbLogFile)
	defer dbHandler.Close() // Defer Closing the database

	if !dbFileExists {
		dbHandler.CreateTable(database.Wifi)                                // Create Database Tables
		dbHandler.InsertRecord(database.Wifi, "some ssid", "some password") // INSERT RECORDS
	}

	dbHandler.DumpTable(database.Wifi) // DISPLAY INSERTED RECORDS
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

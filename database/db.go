package database

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

const (
	Wifi     = 0
	Username = 1
	Email    = 2
)

var WrongKindError = errors.New("wrong record kind")

type DBHandler struct {
	database *sql.DB
	logger   *log.Logger
}

func CreateDBHandler(db *sql.DB, output io.Writer) *DBHandler {
	lg := log.New(output, "INFO: ", log.LstdFlags|log.Lshortfile)
	return &DBHandler{database: db, logger: lg}
}

func (dBhandler *DBHandler) CreateTable(kind int) {
	var tableName string
	var userName string
	switch kind {
	case Wifi:
		tableName = "wifis"
		userName = "SSID"
	case Username:
		tableName = "usernames"
		userName = "username"
	case Email:
		tableName = "emails"
		userName = "emails"
	default:
		dBhandler.logger.Fatal(WrongKindError)
	}

	Query := fmt.Sprintf(`CREATE TABLE %s (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"%s" TEXT,
		"password" TEXT		
	  );`, tableName, userName)

	dBhandler.logger.Println("create wifi table...")

	_, err := dBhandler.database.Exec(Query)
	if err != nil {
		dBhandler.logger.Fatal(err)
	}
	dBhandler.logger.Println("wifi table created")
}

// We are passing db reference connection from main to our method with other parameters
func (dBhandler *DBHandler) InsertRecord(kind int, user string, pass string) {
	var tableName string
	var userName string
	dBhandler.logger.Println("Inserting wifi record ...")
	switch kind {
	case Wifi:
		tableName = "wifis"
		userName = "SSID"
	case Username:
		tableName = "usernames"
		userName = "username"
	case Email:
		tableName = "emails"
		userName = "emails"
	default:
		dBhandler.logger.Fatal(WrongKindError)
	}

	Query := fmt.Sprintf("INSERT INTO %s (%s, password) VALUES (?, ?);", tableName, userName)
	_, err := dBhandler.database.Exec(Query, user, pass)
	if err != nil {
		dBhandler.logger.Fatal(err)
	}
}

func (dBhandler *DBHandler) DumpTable(kind int) {
	var tableName string
	var userName string
	switch kind {
	case Wifi:
		tableName = "wifis"
		userName = "SSID"
	case Username:
		tableName = "usernames"
		userName = "username"
	case Email:
		tableName = "emails"
		userName = "emails"
	default:
		dBhandler.logger.Fatal(WrongKindError)
	}

	Query := fmt.Sprintf("SELECT * FROM %s ORDER BY %s;", tableName, userName)
	row, err := dBhandler.database.Query(Query)
	if err != nil {
		dBhandler.logger.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var ssid string
		var pass string
		row.Scan(&id, &ssid, &pass)
		fmt.Println("Record:\t", id, "\t", ssid, "\t", pass)
	}
}

func (dBhandler *DBHandler) Close() {
	dBhandler.database.Close()
	dBhandler.logger.Println("database closed")
}

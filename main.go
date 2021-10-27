package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/good-times-ahead/password-manager-go/database"
	_ "github.com/lib/pq"
)

// initialize the constants needed for the DB connection
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "passwords"
	table    = "info" // making a constant for the table that we shall be using

)

// Check whether the table to use in further operations has already been created
func tableExists(db *sql.DB) error {
	getTable := fmt.Sprintf("select * from %s", table)

	if _, tableErr := db.Query(getTable); tableErr != nil {
		fmt.Println("First-time execution; creating table...")

		if makeTableErr := database.MakeTable(db); makeTableErr != nil {
			return makeTableErr
		}
		fmt.Println("Everything done. You're good to go.")
	}
	return nil
}

func main() {

	db, connErr := database.ConnectToDB(host, port, user, password, dbname)
	if connErr != nil {
		log.Fatal("Couldn't connect to database!")
	}

	checkTable := tableExists(db)

	if checkTable != nil {
		log.Fatal(checkTable)
	}

}

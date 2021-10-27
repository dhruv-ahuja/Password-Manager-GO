package main

import (
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

func main() {

	db, connErr := database.ConnectToDB(host, port, user, password, dbname)
	if connErr != nil {
		log.Fatal("Couldn't connect to database!")
	}

	checkTable := database.TableExists(db, table)

	if checkTable != nil {
		log.Fatal(checkTable)
	}

}

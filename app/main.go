package main

import (
	"fmt"
	"log"
	"os"

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

	// Check whether the table has already been created
	getTable := fmt.Sprintf("select * from %s;", table)

	_, tableErr := db.Query(getTable)
	if tableErr != nil {
		fmt.Println("First-time execution; creating table...")

		if makeTableErr := database.MakeTable(db); makeTableErr != nil {
			log.Fatal("Couldn't create table!")
		}
	} else {
		fmt.Println("Everything done. You're good to go.")
	}

	// makeTable := setupDB()
	// _, err := db.Exec(string(makeTable))
	// if err != nil {
	// 	panic(err)
	// }

}

func setupDB() []byte {
	//Get the path to your sql script in a os agnostic way.
	path := "../db/setup.sql"
	//Read the content of the file.
	makeTable, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	//Execute the statements in the file using the sql driver.*/
	return makeTable
}

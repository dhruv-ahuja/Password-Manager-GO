package database

import (
	"database/sql"
	"fmt"
	"os"

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

var DB *sql.DB

//Connect to database
func ConnecttoDB() error {

	// Prepare postgres connection parameters
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, user, password, dbname)

	// Establish connection
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	// Ping to confirm whether connection works
	if err = db.Ping(); err != nil {
		return err
	}

	fmt.Println("Connected to the Database successfully!")
	fmt.Println()
	DB = db

	return nil

}

// Check whether the table to use exists or not
func TableExists() error {

	// SQL query to run, not prone to injection attacks since we are just inserting the table name manually
	query := fmt.Sprintf("SELECT * FROM %s", table)
	_, err := DB.Exec(query)
	if err != nil {
		// an error means that the table doesn't exist, we need to call the MakeTable function
		if MakeTableErr := MakeTable(); MakeTableErr != nil {
			return MakeTableErr
		}

	} else {
		fmt.Println("Found existing table. Good to go!")
		fmt.Println()
	}

	return nil

}

// Make the table which we will use for all our operations
func MakeTable() error {

	fmt.Println("First-time execution; creating table...")

	// Path to the relevant SQL file
	path := "./database/setup.sql"
	// Read the file content
	queries, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Convert the slice to a string since the database connector only accepts strings for queries.
	if _, err := DB.Exec(string(queries)); err != nil {
		return err
	}

	fmt.Println("Everything done. You're good to go.")
	fmt.Println()

	return nil

}

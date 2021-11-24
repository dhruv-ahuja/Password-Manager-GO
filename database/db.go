package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

//Connect to database
func ConnecttoDB() error {
	// declare necessary parameters
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Prepare postgres connection parameters
	psqlInfo := fmt.Sprint("host=", host, " port=", port, " user=", user, " password=", password, " dbname=", dbname, " sslmode=disable")

	// Establish connection
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		return errors.New("error establishing connection to postgres, please check your parameters")
	}

	// Ping to confirm whether connection works
	if err = db.Ping(); err != nil {
		return errors.New("unable to successfully ping the database")
	}

	fmt.Println("Connected to the Database successfully!")

	DB = db

	return nil

}

// Check whether the table to use exists or not
func TableExists() error {
	// Returns error if table does not exist.
	query := "SELECT 'public.info'::regclass"

	_, err := DB.Exec(query)

	if err != nil {
		// an error means that the table doesn't exist, we need to call the MakeTable function
		if err := MakeTable(); err != nil {
			return err

		}

	} else {
		// adding new lines to keep the interface clean and readable
		fmt.Printf("Found existing table. Good to go!\n\n")
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
		return errors.New("setup.sql file not found or something was modified")
	}

	// Convert the slice to a string since the database connector only accepts strings for queries.
	if _, err := DB.Exec(string(queries)); err != nil {
		return errors.New("unable to make table 'info'")
	}

	// adding new lines to keep the interface clean and readable
	fmt.Printf("Everything done. You're good to go.\n\n")

	return nil

}

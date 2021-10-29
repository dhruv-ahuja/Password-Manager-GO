package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// Connect to Database
func ConnectToDB(host string, port int, username, password, dbname string) (*sql.DB, error) {

	// Prepare postgres connection parameters
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, username, password, dbname)

	// Establish connection
	db, connErr := sql.Open("postgres", psqlInfo)
	if connErr != nil {
		return nil, connErr
	}

	// Ping to confirm connection
	if pingErr := db.Ping(); pingErr != nil {
		return nil, pingErr
	}

	fmt.Println("Connected to DB successfully!")
	return db, nil

}

// Check whether the table to use in further operations has already been created
func TableExists(db *sql.DB, table string) error {

	// SQL query to run
	getTable := fmt.Sprintf("select * from %s", table)

	if _, tableErr := db.Query(getTable); tableErr != nil {

		// if there's no table, means the user is using the application for the first time
		fmt.Println("First-time execution; creating table...")

		// Call the MakeTable function
		if makeTableErr := MakeTable(db); makeTableErr != nil {
			return makeTableErr
		}
		// Let the user know that the table has been created
		fmt.Println("Everything done. You're good to go.")
		fmt.Println()
	}
	return nil
}

// Make table if it doesn't yet exist
func MakeTable(db *sql.DB) error {

	// Get the path to the SQL file
	path := "./database/setup.sql"
	// Read the file content
	queries, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Convert queries to a string since ReadFile returns a slice of bytes but db.Query only accepts string
	if _, tableErr := db.Exec(string(queries)); tableErr != nil {
		return tableErr
	} else {
		return nil
	}
}

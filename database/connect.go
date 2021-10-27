package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// Make table if it doesn't yet exist
func MakeTable(db *sql.DB) error {
	// Get the path to the SQL file
	filepath := "./setup.sql"
	// Read the file content
	queries, err := os.ReadFile(filepath)
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

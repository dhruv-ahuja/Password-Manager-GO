package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Repo struct {
	DB *sql.DB
}

type Config struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
}

// attempting to implement Data Injection/Inversion *sweat*
func NewDBRepo(db *sql.DB) *Repo {
	return &Repo{DB: db}
}

func GenerateConfig() Config {

	var c Config

	c.host = os.Getenv("HOST")
	c.port = os.Getenv("PORT")
	c.user = os.Getenv("DB_USER")
	c.password = os.Getenv("DB_PASSWORD")
	c.dbname = os.Getenv("DB_NAME")

	return c
}

// Check whether the table to use exists or not
func (conn *Repo) TableExists() error {
	// Returns error if table does not exist.
	query := "SELECT 'public.info'::regclass"

	// _, err := DB.Exec(query)
	_, err := conn.DB.Exec(query)

	if err != nil {
		// an error means that the table doesn't exist, we need to call the MakeTable function
		if err := conn.MakeTable(); err != nil {
			return err

		}

	} else {
		// adding new lines to keep the interface clean and readable
		fmt.Printf("Found existing table. Good to go!\n\n")
	}

	return nil

}

// Make the table which we will use for all our operations
func (conn *Repo) MakeTable() error {

	fmt.Println("First-time execution; creating table...")

	// Path to the relevant SQL file
	path := "./database/setup.sql"

	// Read the file content
	queries, err := os.ReadFile(path)

	if err != nil {
		return errors.New("setup.sql file not found or something was modified")
	}

	// Convert the slice to a string since the database connector only accepts strings for queries.
	if _, err := conn.DB.Exec(string(queries)); err != nil {
		return errors.New("unable to make table 'info'")
	}

	// adding new lines to keep the interface clean and readable
	fmt.Printf("Everything done. You're good to go.\n\n")

	return nil

}

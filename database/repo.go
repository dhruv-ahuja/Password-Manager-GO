package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type DBFuncs interface {
	MakeTable(string) error
}

// attempting to implement Data Injection *sweat*
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

func NewDBRepo(db *sql.DB) *Repo {
	return &Repo{DB: db}
}

// NewConfig returns a new Config instance containing
// all necessary connection parameters
func NewConfig() Config {

	return Config{
		host:     os.Getenv("HOST"),
		port:     os.Getenv("PORT"),
		user:     os.Getenv("DB_USER"),
		password: os.Getenv("DB_PASSWORD"),
		dbname:   os.Getenv("DB_NAME"),
	}
}

//Connect to database, return the connection variable for usage throughout program
func NewConnection(c Config) (*sql.DB, error) {

	// Prepare postgres connection parameters
	psqlInfo := fmt.Sprint("host=", c.host, " port=", c.port,
		" user=", c.user, " password=", c.password,
		" dbname=", c.dbname, " sslmode=disable")

	// Establish connection
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		return nil, errors.New("error connecting to postgres, please check")
	}

	// Ping to confirm whether connection works
	if err = db.Ping(); err != nil {
		return nil, errors.New("unable to ping the database")
	}

	fmt.Println("Connected to the Database successfully!")

	return db, nil

}

// Make the table which we will use for all our operations
func (conn *Repo) MakeTable(sqlFilePath string) error {

	fmt.Println("First-time execution; creating table...")

	// Read the file content
	query, err := os.ReadFile(sqlFilePath)

	if err != nil {
		return errors.New("setup.sql file not found or something was modified")
	}

	// Convert the slice to a string since the database connector only accepts strings for queries.
	if _, err := conn.DB.Exec(string(query)); err != nil {
		return errors.New("unable to make table 'info'")
	}

	// adding new lines to keep the interface clean and readable
	fmt.Printf("Everything done. You're good to go.\n\n")

	return nil

}

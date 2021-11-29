package database

import (
	"database/sql"
	"errors"
	"fmt"
)

//Connect to database
func ConnecttoDB(c Config) (*sql.DB, error) {

	// Prepare postgres connection parameters
	psqlInfo := fmt.Sprint("host=", c.host, " port=", c.port, " user=", c.user, " password=", c.password, " dbname=", c.dbname, " sslmode=disable")

	// Establish connection
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		return nil, errors.New("error establishing connection to postgres, please check your parameters")
	}

	// Ping to confirm whether connection works
	if err = db.Ping(); err != nil {
		return nil, errors.New("unable to ping the database")
	}

	fmt.Println("Connected to the Database successfully!")

	return db, nil

}

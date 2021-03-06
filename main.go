package main

import (
	"fmt"
	"log"

	"github.com/good-times-ahead/password-manager-go/auth"
	"github.com/good-times-ahead/password-manager-go/program"
	"github.com/good-times-ahead/password-manager-go/store"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	// Path to hashed master password file, encrypted data and SQL file
	pwFilePath := "./master_pw"
	encInfoPath := "./encrypted_data"
	sqlFilePath := "./setup.sql"

	if err := godotenv.Load(); err != nil {
		log.Fatal(
			fmt.Errorf("error reading from .env file: %s", err),
		)
	}

	dbStore, err := store.NewDBStore()

	if err != nil {
		log.Fatal(
			fmt.Errorf("error initializing DBStore: %s", err),
		)
	}

	// Get a new struct instance
	cli := program.New(dbStore)

	// Init runs all initial checks and also sets up the database connection through store.DBStore
	err = cli.Init(pwFilePath, encInfoPath, sqlFilePath)

	if err != nil {
		log.Fatal(err)
	}

	// This phase is only run after ensuring that encryption key, table and master password have been generated
	encryptionKey, err := auth.Run(encInfoPath, pwFilePath)

	if err != nil {
		log.Fatal(err)
	}

	if err := cli.Prompt(encryptionKey); err != nil {
		log.Fatal(err)
	}

}

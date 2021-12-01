package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/good-times-ahead/password-manager-go/app"
	"github.com/good-times-ahead/password-manager-go/auth"
	"github.com/good-times-ahead/password-manager-go/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	err := initialize()

	if err != nil {
		log.Fatal(err)
	}

}

func initialize() error {
	// Load .env file, all loaded environment variables henceforth are available to all functions
	if err := godotenv.Load(); err != nil {
		return errors.New("error reading from .env file, please check")
	}

	// Generate database config
	dbConfig := database.GenerateConfig()

	// Initialize connection to the database
	dbVar, err := database.ConnecttoDB(dbConfig)

	if err != nil {
		return err
	}

	// Get the Struct with DB connection and relevant methods
	repo := database.NewDBRepo(dbVar)

	dbConn := app.NewDBConn(repo)

	// Path to hashed master password file
	pwFilePath := "./master_pw"
	// Path to encrypted data (salt, encryption key)
	encInfoPath := "./encrypted_data"

	// Check whether encrypted data already exists
	checkEncData := auth.CheckEncryptedData(encInfoPath)

	if !checkEncData {

		// Drop any existing table and start afresh
		// err := database.MakeTable()
		err := dbConn.Repo.MakeTable()

		if err != nil {
			return err
		}

		// Execute "first-run" functions
		err = auth.FirstRun(encInfoPath, pwFilePath)

		if err != nil {
			return err
		}

	}

	// Check if our table already exists
	// checkTableErr := database.TableExists()
	checkTableErr := dbConn.Repo.TableExists()

	if checkTableErr != nil {
		return checkTableErr
	}

	err = auth.Run(encInfoPath, pwFilePath)

	if err != nil {
		return err
	}

	// load encrypted data to use when dealing with credentials later
	encData, err := auth.LoadEncryptedInfo(encInfoPath)

	if err != nil {
		return err
	}

	// now, to unseal the encryption key
	encryptionKey, err := auth.UnsealEncryptionKey(pwFilePath, encData)

	if err != nil {
		return err
	}

	// Finally, start the app
	appPersist := true

	for appPersist {
		mainMsg := `Hello, what would you like to do?
1. Save a password to the DB
2. View a saved password
3. Edit a saved password
4. Delete a saved password
0: Exit the application: `

		usrInput := app.GetInput(mainMsg)

		err := dbConn.TakeInput(usrInput, encryptionKey)

		if err != nil {
			return err
		}
		// adding new lines to keep the interface clean and readable
		fmt.Println()
	}

	return nil

}

func Test() {
	fmt.Println("OK")
}

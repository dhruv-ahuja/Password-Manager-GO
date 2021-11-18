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

	// Initialize connection to the database
	initDB := database.ConnecttoDB()

	if initDB != nil {
		return initDB
	}

	// Path to hashed master password file
	pwFilePath := "./master_pw"
	// Path to encrypted data (salt, encryption key)
	encInfoPath := "./encrypted_data"

	// Check if master password exists
	checkPassword := auth.CheckMasterPassword(pwFilePath)

	if !checkPassword {
		err := auth.FirstRun(pwFilePath)

		if err != nil {
			return err
		}
	}

	err := auth.Run(pwFilePath)

	if err != nil {
		return err
	}

	// Check if our table already exists
	checkTableErr := database.TableExists()

	if checkTableErr != nil {
		return checkTableErr
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

	fmt.Println(encryptionKey)

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

		executeAppErr := app.TakeInput(usrInput, encryptionKey)

		if executeAppErr != nil {
			return executeAppErr
		}

		fmt.Println()
	}

	return nil

}

func Test() {
	fmt.Println("OK")
}

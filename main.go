package main

import (
	"fmt"
	"log"

	"github.com/good-times-ahead/password-manager-go/app"
	"github.com/good-times-ahead/password-manager-go/auth"
	"github.com/good-times-ahead/password-manager-go/database"
	_ "github.com/lib/pq"
)

func main() {
	err := initialize()

	if err != nil {
		log.Fatal(err)
	}
}

func initialize() error {
	// Initialize connection to the database
	initDB := database.ConnecttoDB()

	if initDB != nil {
		return initDB
	}

	// Path to hashed master password file
	var pwFilePath = "./master_pw"

	// Check if master password exists
	checkPassword := auth.CheckMasterPassword(pwFilePath)

	if checkPassword != nil {
		return checkPassword
	}

	// Ask user for master password
	authUser := auth.AuthorizeUser(pwFilePath)

	if authUser != nil {
		return authUser
	}

	// Check if our table already exists
	checkForTable := database.TableExists()

	if checkForTable != nil {
		return checkForTable
	}

	// Finally, start the app
	appPersist := true

	for appPersist {
		run := app.TakeInput()

		if run != nil {
			return run
		}

		fmt.Println()
	}

	return nil

}

func Test() {
	fmt.Println("OK")
}

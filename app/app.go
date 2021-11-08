package app

import (
	"fmt"
	"os"
)

// Take user input to dictate what function gets executed
func TakeInput() error {

	mainMsg := `Hello, what would you like to do?
1. Save a password to the DB
2. View a saved password
3. Edit a saved password
0: Exit the application: `

	usrInput := GetInput(mainMsg)

	if usrInput == "1" {
		// returning function directly since it's supposed to return an error anyway
		return SaveCredentials()

	}
	if usrInput == "2" {

		askForWebsite := "Enter the website to retrieve accounts for: "
		website := GetInput(askForWebsite)

		return ViewSavedCredentials(website)

	}
	if usrInput == "3" {

		askForWebsite := "Enter the website to edit credentials for: "
		website := GetInput(askForWebsite)

		return EditCredentials(website)

	}
	if usrInput == "0" {

		os.Exit(0)

	} else {

		fmt.Println("Invalid input!")

	}

	return nil

}

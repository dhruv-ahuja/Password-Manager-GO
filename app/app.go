package app

import (
	"fmt"
	"os"
)

// Take user input to dictate what function gets executed
func TakeInput(usrInput string, encryptionKey []byte) error {

	switch usrInput {

	case "1":
		// returning function directly since it's supposed to return an error anyway
		return SaveCredentials(encryptionKey)

	case "2":
		askForWebsite := "Enter the website to retrieve accounts for: "
		website := GetInput(askForWebsite)

		return ViewSavedCredentials(website, encryptionKey)

	case "3":
		askForWebsite := "Enter the website to edit credentials for: "
		website := GetInput(askForWebsite)

		return EditCredentials(website, encryptionKey)

	case "4":
		askForWebsite := "Enter the website to delete credentials for: "
		website := GetInput(askForWebsite)

		return DeleteCredentials(website, encryptionKey)

	case "0":
		os.Exit(0)

	default:
		fmt.Println("Invalid input!")

	}

	return nil

}

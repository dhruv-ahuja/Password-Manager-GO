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
		askForKey := "Enter the key to retrieve accounts for: "
		key := GetInput(askForKey)

		return ViewSavedCredentials(key, encryptionKey)

	case "3":
		askForKey := "Enter the key to edit credentials for: "
		key := GetInput(askForKey)

		return EditCredentials(key, encryptionKey)

	case "4":
		askForKey := "Enter the website to delete credentials for: "
		key := GetInput(askForKey)

		return DeleteCredentials(key, encryptionKey)

	case "0":
		os.Exit(0)

	default:
		fmt.Println("Invalid input!")

	}

	return nil

}

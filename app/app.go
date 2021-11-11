package app

import (
	"fmt"
	"os"
)

// Take user input to dictate what function gets executed
func TakeInput(usrInput string) error {

	switch usrInput {

	case "1":
		// returning function directly since it's supposed to return an error anyway
		return SaveCredentials()

	case "2":
		askForWebsite := "Enter the website to retrieve accounts for: "
		website := GetInput(askForWebsite)

		return ViewSavedCredentials(website)

	case "3":
		askForWebsite := "Enter the website to edit credentials for: "
		website := GetInput(askForWebsite)

		return EditCredentials(website)

	case "4":
		askForWebsite := "Enter the website to delete credentials for: "
		website := GetInput(askForWebsite)

		return DeleteCredentials(website)

	case "0":
		os.Exit(0)

	default:
		fmt.Println("Invalid input!")

	}

	// if usrInput == "1" {
	// 	// returning function directly since it's supposed to return an error anyway
	// 	return SaveCredentials()

	// }

	// if usrInput == "2" {

	// 	askForWebsite := "Enter the website to retrieve accounts for: "
	// 	website := GetInput(askForWebsite)

	// 	return ViewSavedCredentials(website)

	// }

	// if usrInput == "3" {

	// 	askForWebsite := "Enter the website to edit credentials for: "
	// 	website := GetInput(askForWebsite)

	// 	return EditCredentials(website)

	// }

	// if usrInput == "4" {

	// 	askForWebsite := "Enter the website to delete credentials for: "
	// 	website := GetInput(askForWebsite)

	// 	return DeleteCredentials(website)

	// }

	// if usrInput == "0" {

	// 	os.Exit(0)

	// } else {

	// 	fmt.Println("Invalid input!")

	// }

	return nil

}

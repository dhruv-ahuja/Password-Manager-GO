package app

import "fmt"

// Take user input to dictate what function gets executed
func TakeInput() error {

	mainMsg := `Hello, what would you like to do today?
1. Save a password to the DB.
2. View a saved password.
3. Edit a saved password.`
	fmt.Println(mainMsg)

	var usrInput string
	fmt.Scanf("%s", &usrInput)

	if usrInput == "1" {
		// returning function directly since it's supposed to return an error anyway
		return SaveCredentials()
	}

	return nil
}

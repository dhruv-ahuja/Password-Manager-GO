package app

import (
	"fmt"

	"github.com/good-times-ahead/password-manager-go/database"
)

// View credentials for the specified website
func ViewSavedCredentials(website string) (int, error) {
	// ask user for website name
	// askForWebsite := "Enter the website to retrieve accounts for: "

	// website := GetInput(askForWebsite)

	// Get all accounts associated with the website
	query := "SELECT * FROM info where website=$1"

	rows, err := database.DB.Query(query, website)

	if err != nil {
		return 0, err
	}

	// Prepare a slice to store the retrieved credentials
	var accountsList []credentials

	for rows.Next() {
		var usrInfo credentials
		var base64Password string

		// Write scanned values to credentials except the password, we need to decrypt it first
		if err := rows.Scan(&usrInfo.ID, &usrInfo.website, &usrInfo.email, &usrInfo.username, &base64Password); err != nil {
			return 0, err
		}

		// Now, to decrypt the password
		password, err := usrInfo.DecryptPassword(base64Password)

		if err != nil {
			return 0, err
		}

		// Finally, we write the decrypted password to the credentials struct
		usrInfo.password = password

		// Append the credentials variable to accountList
		accountsList = append(accountsList, usrInfo)

	}

	count := len(accountsList)

	// Print out the results of the query
	if count == 0 {

		fmt.Println("Sorry, no accounts saved for that website!")

	} else {

		for _, usrInfo := range accountsList {
			response := fmt.Sprintf("S.No.: %d, Website: %s, Email: %s, Username: %s, Password: %s", usrInfo.ID, usrInfo.website, usrInfo.email, usrInfo.username, usrInfo.password)

			fmt.Println(response)

		}
	}
	return count, nil

}

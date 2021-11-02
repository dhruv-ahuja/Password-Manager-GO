package app

import (
	"fmt"

	"github.com/good-times-ahead/password-manager-go/database"
)

// View credentials for the specified website
func ViewSavedCredentials(website string) error {
	// Get all accounts associated with the website
	query := "SELECT * FROM info where website=$1"

	rows, err := database.DB.Query(query, website)

	if err != nil {
		return err
	}

	// Prepare a slice to store retrieved credentials
	var accountsList []credentials

	for rows.Next() {

		var usrInfo credentials
		var base64Password string

		// Write scanned values to credentials except the password, we need to decrypt it first
		err := rows.Scan(&usrInfo.ID, &usrInfo.website, &usrInfo.email, &usrInfo.username, &base64Password)

		if err != nil {
			return err
		}

		// Now, to decrypt the password
		password, err := usrInfo.DecryptPassword(base64Password)

		if err != nil {
			return err
		}

		// Finally, we write the decrypted password to the credentials struct
		usrInfo.password = password

		// Append the credentials variable to accountList
		accountsList = append(accountsList, usrInfo)

	}

	// Print out the results of the query
	if len(accountsList) == 0 {

		fmt.Println("Sorry, no accounts saved for that website!")

	} else {

		for _, usrInfo := range accountsList {

			response1 := fmt.Sprintf("ID No. %d, Website: %s", usrInfo.ID, usrInfo.email)

			response2 := fmt.Sprintf("Email: %s, Username: %s, Password: %s", usrInfo.email, usrInfo.username, usrInfo.password)

			fmt.Println(response1 + response2)

		}
	}
	return nil

}

package app

import (
	"github.com/good-times-ahead/password-manager-go/database"
)

// Helper function to retrieve specific data from the table
func RetrieveCredentials(query, website string) ([]credentials, error) {

	rows, err := database.DB.Query(query, website)

	if err != nil {
		return nil, err
	}

	// Prepare a slice to store retrieved credentials
	var accountsList []credentials

	for rows.Next() {

		var usrInfo credentials
		var base64Password string

		// Write scanned values to credentials struct except for the password,
		// which needs to be decrypted
		err := rows.Scan(&usrInfo.ID, &usrInfo.website, &usrInfo.email, &usrInfo.username, &base64Password)

		if err != nil {
			return nil, err
		}

		// Now, to decrypt the password
		password, err := usrInfo.DecryptPassword(base64Password)

		if err != nil {
			return nil, err
		}

		// Finally, we write the decrypted password to its respective field
		usrInfo.password = password

		// Append the credentials variable to the slice
		accountsList = append(accountsList, usrInfo)
	}

	return accountsList, nil

}

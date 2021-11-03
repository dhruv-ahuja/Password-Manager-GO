package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/good-times-ahead/password-manager-go/database"
)

// Function to get user input in a streamlined fashion.
func GetInput(argument string) string {

	// Emulate a while loop to receive user input and ensure its' validity
	reader := bufio.NewReader(os.Stdin)

	isEmpty := true

	for isEmpty {

		fmt.Print(argument)

		usrInput, err := reader.ReadString('\n')
		//fmt.Println(usrInput, len(usrInput))

		// We want to keep re-iterating over the loop so we can leave the error as is(I think)
		if err != nil {
			fmt.Println("Invalid input or method!")
		}

		switch len(usrInput) {

		case 1:
			fmt.Println("Empty input!")

		default:
			return strings.TrimSpace(usrInput)
		}

	}

	return ""

}

// Function to retrieve specific data from the table
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

	if len(accountsList) == 0 {
		// if the received slice is empty
		fmt.Println("Sorry, no accounts saved for that website!")

	} else {

		PrintEntries(accountsList)

	}

	return accountsList, nil

}

// Print entries received from database queries
func PrintEntries(accountsList []credentials) {

	// print out the list of found entries
	for _, usrInfo := range accountsList {
		// dividing response string into 2 parts to maintain visibility
		response1 := fmt.Sprintf("ID No. %d, Website: %s, ", usrInfo.ID, usrInfo.website)

		response2 := fmt.Sprintf("Email: %s, Username: %s, Password: %s", usrInfo.email, usrInfo.username, usrInfo.password)

		fmt.Println(response1 + response2)

	}
}

package app

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/good-times-ahead/password-manager-go/database"
	"golang.org/x/term"
)

// Function to get user input in a streamlined fashion.
func GetInput(argument string) string {

	// Emulate a while loop to receive user input and ensure its' validity
	reader := bufio.NewReader(os.Stdin)

	isEmpty := true

	for isEmpty {

		fmt.Print(argument)

		usrInput, err := reader.ReadString('\n')

		// We want to keep re-iterating over the loop so we can leave the error as is(I think)
		if err != nil {
			fmt.Println("Invalid input or method!")
		}

		switch len(usrInput) {

		case 1:
			fmt.Println("Empty input!")
			fmt.Println()
		default:
			fmt.Println()
			return strings.TrimSpace(usrInput)

		}

	}

	return ""

}

// GetPassInput takes input using term.ReadPassword(), preventing user input echo
func GetPassInput(argument string) []byte {

	// Emulate a while loop to receive user input and ensure its' validity
	isEmpty := true

	for isEmpty {

		fmt.Print(argument)

		usrInput, err := term.ReadPassword(int(syscall.Stdin))

		if err != nil {
			fmt.Println("Invalid input or method!")
		}

		switch len(usrInput) {

		case 0:
			fmt.Println("Empty input!")
			fmt.Println()
		default:
			return usrInput

		}

	}

	return nil

}

// Function to retrieve specific data from the table
func RetrieveCredentials(query, website string, encryptionKey []byte) ([]Credentials, error) {

	rows, err := database.DB.Query(query, website)

	if err != nil {
		return nil, errors.New("error executing query")
	}

	// Prepare a slice to store retrieved credentials
	var accountsList []Credentials

	for rows.Next() {

		var usrInfo Credentials
		var base64Password string

		// Write scanned values to credentials struct except for the password,
		// which needs to be decrypted
		err := rows.Scan(&usrInfo.ID, &usrInfo.website, &usrInfo.email, &usrInfo.username, &base64Password)

		if err != nil {
			return nil, errors.New("error attempting to retrieve data from query")
		}

		// Now, to decrypt the password
		password, err := usrInfo.DecryptPassword(base64Password, encryptionKey)

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
func PrintEntries(accountsList []Credentials) {

	// print out the list of found entries
	for _, usrInfo := range accountsList {
		// dividing response string into 2 parts to maintain visibility
		response1 := fmt.Sprintf("ID No. %d, Website: %s, ", usrInfo.ID, usrInfo.website)

		response2 := fmt.Sprintf("Email: %s, Username: %s, Password: %s", usrInfo.email, usrInfo.username, usrInfo.password)

		fmt.Println(response1 + response2)
		fmt.Println()

	}
}

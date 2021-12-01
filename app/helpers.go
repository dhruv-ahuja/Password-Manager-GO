package app

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/good-times-ahead/password-manager-go/credentials"
	"github.com/good-times-ahead/password-manager-go/password"
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
			fmt.Printf("Empty input!\n\n")
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
			fmt.Printf("Empty input!\n\n")
		default:
			return usrInput

		}

	}

	return nil

}

// Function to retrieve specific data from the table
func (DBConn *DBConn) RetrieveCredentials(query, key string, encryptionKey []byte) ([]credentials.Credentials, error) {

	rows, err := DBConn.Repo.DB.Query(query, "%"+key+"%")

	if err != nil {
		return nil, errors.New("error executing query")
	}

	// Prepare a slice to store retrieved credentials
	var accountsList []credentials.Credentials

	for rows.Next() {

		var usrInfo credentials.Credentials
		var base64Password string

		// Write scanned values to credentials struct except for the password,
		// which needs to be decrypted
		err := rows.Scan(&usrInfo.ID, &usrInfo.Key, &base64Password)

		if err != nil {
			return nil, errors.New("error attempting to retrieve data from query")
		}

		// Now, to decrypt the password
		password, err := password.Decrypt(base64Password, encryptionKey, usrInfo)

		if err != nil {
			return nil, err
		}

		// Finally, we write the decrypted password to its respective field
		usrInfo.Password = password

		// Append the credentials variable to the slice
		accountsList = append(accountsList, usrInfo)
	}

	if len(accountsList) == 0 {
		// if the received slice is empty
		fmt.Println("Sorry, no accounts saved for that key!")

	} else {

		PrintEntries(accountsList)

	}

	return accountsList, nil

}

// Print entries received from database queries
func PrintEntries(accountsList []credentials.Credentials) {

	// print out the list of found entries
	for _, usrInfo := range accountsList {
		// dividing response string into 2 parts to maintain visibility
		response1 := fmt.Sprintf("ID No. %d, Key: %s, Password: %s", usrInfo.ID, usrInfo.Key, usrInfo.Password)

		fmt.Println(response1)
		fmt.Println()

	}
}

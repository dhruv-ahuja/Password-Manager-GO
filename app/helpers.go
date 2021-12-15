package app

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/good-times-ahead/password-manager-go/password"
	"golang.org/x/term"
)

// Function to retrieve specific data from the table
func (p *Program) RetrieveCredentials(query, key string, encryptionKey []byte) ([]map[string]string, error) {

	rows, err := p.Repo.DB.Query(query, "%"+key+"%")

	if err != nil {
		return nil, errors.New("error executing query")
	}

	// Prepare a slice of maps to store retrieved credentials
	var accountsList []map[string]string

	for rows.Next() {

		usrInfo := make(map[string]string, 3)
		var id, key, base64Password string

		// Write scanned values to map
		err := rows.Scan(&id, &key, &base64Password)

		if err != nil {
			return nil, fmt.Errorf("error attempting to retrieve data from query: %s", err)
		}

		// Now, to decrypt the password
		password, err := password.Decrypt(base64Password, encryptionKey)

		if err != nil {
			return nil, err
		}

		// Finally, we write the decrypted password plus other variables to their respective keys
		usrInfo["id"] = id
		usrInfo["key"] = key
		usrInfo["password"] = password

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
func PrintEntries(accountsList []map[string]string) {

	// print out the list of found entries (in the form of dictionaries/maps)
	for _, usrInfo := range accountsList {

		response := fmt.Sprintf("ID No. %s, Key: %s, Password: %s", usrInfo["id"], usrInfo["key"], usrInfo["password"])

		fmt.Println(response)
		fmt.Println()

	}
}

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

// GetPassInput receives user's input securely, without any echo
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

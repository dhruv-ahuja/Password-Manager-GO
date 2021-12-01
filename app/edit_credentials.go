package app

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/good-times-ahead/password-manager-go/credentials"
	"github.com/good-times-ahead/password-manager-go/password"
)

// Function to allow the user to edit credentials
func (DBConn *DBConn) EditCredentials(key string, encryptionKey []byte) error {

	query := "SELECT * FROM info WHERE key ILIKE $1 ORDER BY id ASC;"

	// call the function to retrieve credentials given relevant query
	accountsList, err := DBConn.RetrieveCredentials(query, key, encryptionKey)

	if err != nil {
		return err
	}

	// if accountsList is empty, return nil since there are no entries to modify
	if len(accountsList) == 0 {
		return nil
	}

	selectID := false
	input := 0
	for !selectID {

		// Get users' input to find the entry they want to modify
		msg := "Enter the ID No. of the entry you want to modify: "

		usrInput := GetInput(msg)

		// Converting the string input to integer for comparison
		input, err = strconv.Atoi(usrInput)

		if err != nil {
			return errors.New("error converting user input(string) to integer")
		}

		for _, entry := range accountsList {
			if entry.ID == input {
				selectID = true
				break
			}
		}
		fmt.Println("Entered ID outside range!")
	}

	// TODO: add check to ensure id in accountsList
	fmt.Println("Input is: ", input)
	// Preparing struct variable to store users' desired entry
	var selection credentials.Credentials

	for _, usrInfo := range accountsList {
		if input == usrInfo.ID {

			selection = usrInfo

		}
	}

	// Now, we have the users' choice of entry, allow them to edit whatever field they want
	// Print out each field alongside the current counter-part
	// Using bufio NewReader since GetInput function doesn't accept empty input
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Your current 'key' entry is: ", selection.Key)
	fmt.Print("Enter new key (leave field blank if no changes): ")

	newKey, err := reader.ReadString('\n')

	if err != nil {
		return err
	}

	fmt.Println("Your current password is: ", selection.Password)
	passPrompt := "Enter new password: "

	// newPassword, err := reader.ReadString('\n')
	getUserPass := GetPassInput(passPrompt)
	newPassword := string(getUserPass)

	// Trim away spaces left behind by user and ReadString function
	newKey = strings.TrimSpace(newKey)

	// if no errors occur, update current values and prepare to update database entry
	// todo: try and implement a better way of doing this
	if newKey != "" {
		selection.Key = newKey
	}

	var b64Password string

	if newPassword != "" {

		selection.Password = newPassword

		// // encrypt updated password
		// b64Password, err = selection.EncryptPassword(encryptionKey)
		b64Password, err = password.Encrypt(encryptionKey, selection)

		if err != nil {
			return err
		}
	}

	// Now, to finally save the new details

	// declaring flag modifyPassword
	modifyPassword := false

	if b64Password != "" {
		// Set b64Password as struct field only if modified (i.e., not empty)
		selection.Password = b64Password

		// set modifyPassword to true
		modifyPassword = true
		err = selection.UpdateCredentials(modifyPassword)

	} else {
		err = selection.UpdateCredentials(modifyPassword)
	}

	if err != nil {
		return err
	}

	fmt.Println("Updated your credentials successfully!")

	return nil
}

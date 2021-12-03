package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/good-times-ahead/password-manager-go/password"
)

// Function to allow the user to edit credentials
func (p *Program) EditCredentials(key string, encryptionKey []byte) error {

	query := "SELECT * FROM info WHERE key ILIKE $1 ORDER BY id ASC;"

	// call the function to retrieve credentials given relevant query
	accountsList, err := p.RetrieveCredentials(query, key, encryptionKey)

	if err != nil {
		return err
	}

	// if accountsList is empty, return nil since there are no entries to modify
	if len(accountsList) == 0 {
		return nil
	}

	selectID := false
	usrInput := ""

	for !selectID {

		// Get users' input to find the entry they want to modify
		msg := "Enter the ID No. of the entry you want to modify: "

		usrInput = GetInput(msg)

		for _, entry := range accountsList {

			if entry["id"] == usrInput {
				selectID = true
			}
		}

		if selectID {
			break
		}

		fmt.Println("Entered ID outside range!")
	}

	// Preparing struct variable to store users' desired entry
	// var selection credentials.Credentials
	selection := make(map[string]string, 3)

	for _, usrInfo := range accountsList {
		if usrInput == usrInfo["id"] {

			selection = usrInfo

		}
	}

	// Now, we have the users' choice of entry, allow them to edit whatever field they want
	// Print out each field alongside the current counter-part
	// Using bufio NewReader since GetInput function doesn't accept empty input
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Your current 'key' entry is: ", selection["key"])
	fmt.Print("Enter new key (leave field blank if no changes): ")

	newKey, err := reader.ReadString('\n')

	if err != nil {
		return err
	}

	fmt.Println("Your current password is: ", selection["password"])
	passPrompt := "Enter new password: "

	getUserPass := GetPassInput(passPrompt)
	newPassword := string(getUserPass)

	// Trim away spaces left behind by user and ReadString function
	newKey = strings.TrimSpace(newKey)

	// Update current values and prepare to update database entry
	if newKey != "" {
		selection["key"] = newKey
	}

	var b64Password string

	if newPassword != "" {

		selection["password"] = newPassword

		// encrypt updated password
		b64Password, err = password.Encrypt(encryptionKey, selection["password"])

		if err != nil {
			return err
		}
	}

	// Now, to finally save the new details

	// declaring flag modifyPassword
	modifyPassword := false

	if b64Password != "" {
		// Set b64Password as struct field only if modified (i.e., not empty)
		selection["password"] = b64Password

		// set modifyPassword to true
		modifyPassword = true
		// err = UpdateCredentials(modifyPassword)
		err = p.Repo.UpdateCredentials(modifyPassword, selection)

	} else {
		err = p.Repo.UpdateCredentials(modifyPassword, selection)
	}

	if err != nil {
		return err
	}

	fmt.Println("Updated your credentials successfully!")

	return nil
}

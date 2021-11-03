package app

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// Function to allow the user to edit credentials
func EditCredentials(website string) error {

	query := "SELECT * FROM info where website=$1"

	// call the function to retrieve credentials given relevant query
	accountsList, err := RetrieveCredentials(query, website)

	if err != nil {
		return err
	}

	// Get users' input to find the entry they want to modify
	msg := ("Enter the ID No. of the entry you want to modify: ")

	usrInput := GetInput(msg)

	// Converting the string input to integer for comparison
	input, err := strconv.Atoi(usrInput)

	if err != nil {
		return err
	}

	// Preparing struct variable to store users' desired entry
	var selection credentials

	for _, usrInfo := range accountsList {
		if input == usrInfo.ID {

			selection = usrInfo

		}
	}

	// Now, we have the users' choice of entry, allow them to edit whatever field they want
	// Print out each field alongside the current counter-part
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Your current email is: ", selection.email)
	fmt.Print("Enter new email (leave field blank if no changes): ")

	newEmail, err := reader.ReadString('\n')

	if err != nil {
		return err
	}

	fmt.Println("Your current username is: ", selection.username)
	fmt.Print("Enter new username (leave field blank if no changes): ")

	newUsername, err := reader.ReadString('\n')

	if err != nil {
		return err
	}

	fmt.Println("Your current password is: ", selection.password)
	fmt.Print("Enter new password (leave field blank if no changes): ")

	newPassword, err := reader.ReadString('\n')

	if err != nil {
		return err
	}

	// if no errors occur, update current values and prepare to update database entry

	fmt.Print(newEmail, newPassword, newUsername)
	fmt.Println()

	return nil
}

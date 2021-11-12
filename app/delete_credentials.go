package app

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/good-times-ahead/password-manager-go/database"
)

// Function to allow users to delete credentials
func DeleteCredentials(website string) error {
	//"$" is postgres' equivalent of "?"
	query := "SELECT * FROM info where website=$1 ORDER BY id ASC"

	// call the function to retrieve credentials given relevant query
	accountsList, err := RetrieveCredentials(query, website)

	if err != nil {
		return err
	}

	// if accountsList is empty, return nil since there are no entries to modify
	if len(accountsList) == 0 {
		return nil
	}

	// Get users' input to find the entry they want to delete
	msg := "Enter the ID No. of the entry you want to delete: "

	usrInput := GetInput(msg)

	// Converting the string input to integer for comparison
	input, err := strconv.Atoi(usrInput)

	if err != nil {
		return errors.New("error converting user input(string) to integer")
	}

	deletionQuery := "DELETE FROM info WHERE ID=$1"

	_, err = database.DB.Exec(deletionQuery, input)

	if err != nil {
		return errors.New("error deleting selected entry")
	}

	fmt.Println("Successfully deleted selected entry!")

	return nil

}

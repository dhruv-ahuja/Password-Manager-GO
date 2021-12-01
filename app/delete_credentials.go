package app

import (
	"errors"
	"fmt"
)

// Function to allow users to delete credentials
func (DBConn *DBConn) DeleteCredentials(key string, encryptionKey []byte) error {

	//"$" is postgres' equivalent of "?"
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

		// Get users' input to find the entry they want to delete
		msg := "Enter the ID No. of the entry you want to delete: "

		usrInput := GetInput(msg)

		// Converting the string input to integer for comparison
		// input, err = strconv.Atoi(usrInput)

		if err != nil {
			return errors.New("error converting user input(string) to integer")
		}

		for _, entry := range accountsList {
			if entry["id"] == usrInput {
				selectID = true
				break
			}
		}
		fmt.Println("Entered ID outside range!")
	}

	deletionQuery := "DELETE FROM info WHERE ID=$1"

	_, err = DBConn.Repo.DB.Exec(deletionQuery, input)

	if err != nil {
		return errors.New("error deleting selected entry")
	}

	fmt.Println("Successfully deleted selected entry!")

	return nil

}

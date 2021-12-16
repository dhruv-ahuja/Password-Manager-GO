package app

import (
	"errors"
	"fmt"
)

// Function to allow users to delete credentials
func (p *Program) DeleteCredentials(key string, encryptionKey []byte) error {

	//"$" is postgres' equivalent of "?"
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

		// Get users' input to find the entry they want to delete
		msg := "Enter the ID No. of the entry you want to delete: "

		usrInput = GetInput(msg)

		for _, entry := range accountsList {
			fmt.Println(entry)
			if entry["id"] == usrInput {
				selectID = true
			}
		}

		if selectID {
			break
		}

		fmt.Println("Entered ID outside range!")
	}

	deletionQuery := "DELETE FROM info WHERE ID=$1"

	_, err = p.Repo.DB.Exec(deletionQuery, usrInput)

	if err != nil {
		return errors.New("error deleting selected entry")
	}

	fmt.Println("Successfully deleted selected entry!")

	return nil

}

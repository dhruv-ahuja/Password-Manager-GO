package program

import (
	"errors"
	"fmt"
	"os"

	"github.com/good-times-ahead/password-manager-go/auth"
	"github.com/good-times-ahead/password-manager-go/store"
)

type Program struct {
	store store.Store
}

func New() Program {
	return Program{}
}

// Init injects the store interface into Program
func (p *Program) Init(pwFilePath, encInfoPath, sqlFilePath string) error {

	dbStore := store.NewDBStore()

	// connecting to postgres database
	err := dbStore.NewConn()

	if err != nil {
		return err
	}

	p.store = &dbStore

	// running necessary checks and the like
	checkEncData := auth.CheckEncryptedData(encInfoPath)

	if !checkEncData {

		if err := dbStore.CreateTable(sqlFilePath); err != nil {
			return err
		}

		if err := auth.FirstRun(encInfoPath, pwFilePath); err != nil {
			return err
		}

	}

	return nil

}

func (p *Program) Prompt(encryptionKey []byte) error {

	if p.store == nil {
		return errors.New("Store not initialized")
	}

	appPersist := true

	for appPersist {

		for appPersist {
			mainMsg := `Hello, what would you like to do?
1. Save a password to the DB
2. View a saved password
3. Edit a saved password
4. Delete a saved password
0: Exit the application: `

			usrInput := store.GetInput(mainMsg)

			if err := p.controller(usrInput, encryptionKey); err != nil {
				return err
			}
			// adding new lines to keep the interface clean and readable
			fmt.Println()

		}

	}

	return nil
}

func (p *Program) controller(usrInput string, encryptionKey []byte) error {

	switch usrInput {

	case "1":
		return p.store.SaveCreds(encryptionKey)

	case "2":
		askForKey := "Enter the key to retrieve accounts for: "
		key := store.GetInput(askForKey)

		return p.store.ViewCreds(key, encryptionKey)

	case "3":
		askForKey := "Enter the key to edit credentials for: "
		key := store.GetInput(askForKey)

		return p.store.EditCreds(key, encryptionKey)

	case "4":
		askForKey := "Enter the website to delete credentials for: "
		key := store.GetInput(askForKey)

		return p.store.DeleteCreds(key, encryptionKey)

	case "0":
		os.Exit(0)

	default:
		fmt.Println("Invalid input!")

	}

	return nil

}

package app

import (
	"fmt"
	"os"

	"github.com/good-times-ahead/password-manager-go/database"
)

type AppFuncs interface {
	TakeInput(string, [][]byte) error

	SaveCredentials([]byte) error

	ViewCredentials(string, []byte) error

	EditCredentials(string, []byte) error

	DeleteCredentials(string, []byte) error
}

// struct that stores a struct that stores the database connection
type DBConn struct {
	Repo database.Repo
}

func NewDB(db database.Repo) *DBConn {
	return &DBConn{
		Repo: db,
	}
}

// Take user input to dictate what function gets executed
func (DBConn *DBConn) TakeInput(usrInput string, encryptionKey []byte) error {

	switch usrInput {

	case "1":
		// returning function directly since it's supposed to return an error anyway
		return DBConn.SaveCredentials(encryptionKey)

	case "2":
		askForKey := "Enter the key to retrieve accounts for: "
		key := GetInput(askForKey)

		return DBConn.ViewCredentials(key, encryptionKey)

	case "3":
		askForKey := "Enter the key to edit credentials for: "
		key := GetInput(askForKey)

		return DBConn.EditCredentials(key, encryptionKey)

	case "4":
		askForKey := "Enter the website to delete credentials for: "
		key := GetInput(askForKey)

		return DBConn.DeleteCredentials(key, encryptionKey)

	case "0":
		os.Exit(0)

	default:
		fmt.Println("Invalid input!")
		DBConn.test()

	}

	return nil

}

func (DB *DBConn) test() {
	fmt.Println("OK")
}

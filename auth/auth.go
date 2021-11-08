package auth

import (
	"errors"
	"fmt"
	"os"

	"github.com/good-times-ahead/password-manager-go/app"
	"github.com/good-times-ahead/password-manager-go/database"
	"golang.org/x/crypto/bcrypt"
)

// Check whether master password file exists already
func CheckMasterPassword() error {
	// Path to hashed master password file
	path := "./master_pw"

	// Check for file
	_, err := os.Open(path)

	if err != nil {
		// Call the function to make master password
		err := MakeMasterPassword(path)

		if err != nil {
			return errors.New("unable to create master password")
		}

		// Also create the table to use
		err = database.MakeTable()

		if err != nil {
			return err
		}

	}
	// Return nil if all is well
	return nil

}

func MakeMasterPassword(path string) error {

	msg := `
Hello and welcome to the Password Manager GO application. If you are seeing this message then this must be your first time using the application. 
To get started, you must first create a master password which will be used to authenticate you each time you run the application.
Set a secure password and remember it since there will be no way to recover it!
	`

	// Print out the introductory message
	fmt.Println(msg)

	// use App packages' GetInput function to get user input
	prompt := "Enter the Master Password: "

	usrInput := app.GetInput(prompt)

	// Generate hashed password from the received user input
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usrInput), 13)

	if err != nil {
		return err
	}

	// 0777 perms allow read write & execute for owner, groups and others
	err = os.WriteFile(path, hashedPassword, 0777)

	if err != nil {
		return err
	}

	fmt.Println("Successfully saved master password to file! Now you will be asked to enter it each time you run the program.")
	return nil

}

// Take user's master password and compare it to the stored hash, allowing or disallowing them access to the application
func AuthorizeUser() error {
	flag := true
	for flag {
		// Take user input
		prompt := "Enter the Master Password: "
		usrInput := app.GetInput(prompt)

		// Load the hash from the file
		path := "./master_pw"
		hashedPassword, err := os.ReadFile(path)

		if err != nil {
			return errors.New("error encountered when attempting to read from the password file")
		}

		// Compare hash and password
		compare := bcrypt.CompareHashAndPassword(hashedPassword, []byte(usrInput))

		switch compare {
		case nil:
			flag = false

		default:
			fmt.Println("The passwords do not match! Try again.")
		}
	}

	return nil

}

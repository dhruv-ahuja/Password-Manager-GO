package auth

import (
	"errors"
	"fmt"
	"os"

	"github.com/good-times-ahead/password-manager-go/app"
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
		makePassword := MakeMasterPassword(path)

		if makePassword != nil {
			return errors.New("unable to create master password")
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usrInput), 5)

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

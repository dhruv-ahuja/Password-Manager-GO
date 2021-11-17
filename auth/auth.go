package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"

	"github.com/good-times-ahead/password-manager-go/app"
	"github.com/good-times-ahead/password-manager-go/database"
	"golang.org/x/crypto/bcrypt"
)

// Check whether master password file exists already
func CheckMasterPassword(pwFilePath string) error {
	// Check for file
	_, err := os.Open(pwFilePath)

	if err != nil {
		// Call the function to make master password
		err := MakeMasterPassword(pwFilePath)

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

func MakeMasterPassword(pwFilePath string) error {

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
	err = os.WriteFile(pwFilePath, hashedPassword, 0777)

	if err != nil {
		return err
	}

	fmt.Println("Successfully saved master password to file! Now you will be asked to enter it each time you run the program.")
	return nil

}

// Take user's master password and compare it to the stored hash, allowing or disallowing them access to the application
func AuthorizeUser(pwFilePath string) error {
	// Load the hash from the file
	hashedPassword, err := os.ReadFile(pwFilePath)

	if err != nil {
		return errors.New("error encountered when attempting to read from the password file")
	}

	// Infinite loop till the user correctly enters the password
	for true {
		// Take user input
		prompt := "Enter the Master Password: "
		usrInput := app.GetInput(prompt)

		// Compare hash and password, returns nil if match else error
		compare := bcrypt.CompareHashAndPassword(hashedPassword, []byte(usrInput))

		if compare == nil {
			return nil
		}

		fmt.Println("The passwords do not match! Try again.")
		fmt.Println()

	}

	return nil

}

// Intended to only be used on first run
func GenerateEncryptionKey() error {
	// add file path to check for existing enc_key existence

	// if file empty/doesn't exist, generate a new encryption key
	encryptionKey := make([]byte, 32)

	if _, err := rand.Read(encryptionKey); err != nil {
		return err
	}
	fmt.Println(encryptionKey)
	return nil
}

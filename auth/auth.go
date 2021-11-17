package auth

import (
	"errors"
	"fmt"
	"os"

	"github.com/good-times-ahead/password-manager-go/app"
	"golang.org/x/crypto/bcrypt"
)

// Check whether master password file exists already
func CheckMasterPassword(pwFilePath string) bool {
	// Check for file's existence, returns an error if unable to open
	checkFile, err := os.OpenFile(pwFilePath, os.O_RDONLY, 0777)
	if err != nil {
		return false
	}

	// Read file to confirm there are 32 bytes of the hashed master password
	isEmpty, err := checkFile.Read(make([]byte, 32))

	// If not, return false since either the file has been tampered with
	if err != nil || isEmpty == 0 {
		return false
	}

	return true

}

// Take master password hash and compare it to the stored hash, allowing or disallowing them access to the application
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

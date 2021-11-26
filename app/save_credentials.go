package app

import (
	"github.com/good-times-ahead/password-manager-go/credentials"
	"github.com/good-times-ahead/password-manager-go/password"
)

// Save credentials to the database
func SaveCredentials(encryptionKey []byte) error {

	// define needed prompts
	promptWebsite := "Enter the websites' name: "
	promptEmail := "Enter your mail ID: "
	promptUsername := "Enter your username: "
	promptPassword := "Enter your password(it will be encrypted before saving): "

	// initialize the variable to save the credentials to
	var usrInfo credentials.Credentials

	// Write user input to respective structure fields
	usrInfo.Website = GetInput(promptWebsite)
	usrInfo.Email = GetInput(promptEmail)
	usrInfo.Username = GetInput(promptUsername)
	usrInfo.Password = string(GetPassInput(promptPassword))

	// encrypt the plain text password
	// encryptedPassword, err := usrInfo.EncryptPassword(encryptionKey)
	encryptedPassword, err := password.Encrypt(encryptionKey, usrInfo)

	if err != nil {
		return err
	}

	// save the credentials to the database
	if saveToDB := usrInfo.InsertIntoDB(encryptedPassword); saveToDB != nil {
		return saveToDB
	}

	return nil

}

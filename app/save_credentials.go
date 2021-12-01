package app

import (
	"github.com/good-times-ahead/password-manager-go/password"
)

// Save credentials to the database
func (DBConn *DBConn) SaveCredentials(encryptionKey []byte) error {

	// define needed prompts
	promptKey := "Enter key: "
	promptPassword := "Enter your password(it will be encrypted before saving): "

	// initialize the variable to save the credentials to
	usrInfo := make(map[string]string, 3)

	// Write user input to respective structure fields
	usrInfo["key"] = GetInput(promptKey)
	usrInfo["password"] = string(GetPassInput(promptPassword))

	// encrypt the plain text password
	encryptedPassword, err := password.Encrypt(encryptionKey, usrInfo["password"])

	if err != nil {
		return err
	}

	// save the credentials to the database
	if saveToDB := DBConn.Repo.InsertIntoDB(encryptedPassword, usrInfo); saveToDB != nil {
		return saveToDB
	}

	return nil

}

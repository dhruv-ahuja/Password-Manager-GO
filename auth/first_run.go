package auth

import (
	"crypto/rand"
	"fmt"
	"os"

	"github.com/good-times-ahead/password-manager-go/app"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/nacl/secretbox"
)

// This file contains functions that will be executed if it is the users' run of the app

func GetMasterPassword() string {

	msg := `
Hello and welcome to the Password Manager GO application. If you are seeing this message then this must be your first time using the application. 
To get started, you must first create a master password which will be used to authenticate you each time you run the application.
Set a secure password and remember it since there will be no way to recover it!
	`

	// Print out the introductory message
	fmt.Println(msg)

	// get users' desired master password in plain text, will be hashed later
	// use App packages' GetInput function
	prompt := "Enter the Master Password: "

	usrInput := app.GetInput(prompt)

	return usrInput

}

// Argon2 is considered better than bcrypt for securing passwords
func HashMasterPassword(masterPassword, pwFilePath string) ([]byte, []byte, error) {
	// generate 32 bytes salt
	salt := make([]byte, 32)

	// generate 32 byte long random salt
	// rand.Read writes the output to the given slice
	if _, err := rand.Read(salt); err != nil {
		return nil, nil, err
	}

	// generate hashed password using argon2
	hashedMasterPassword := argon2.IDKey([]byte(masterPassword), salt, 1, 64*1024, 4, 32)

	// returning salt as well as the hashed master password.
	// the salt will be written to disk as well
	return salt, hashedMasterPassword, nil
}

func SaveMasterPassword(pwFilePath string, hashedMasterPassword []byte) error {
	// Creates file if doesn't exist; permission code "4" means the file is read-only
	err := os.WriteFile(pwFilePath, hashedMasterPassword, 0444)

	if err != nil {
		return err
	}

	fmt.Println("Successfully saved master password to file! Now you will be asked to enter it each time you run the program.")
	return nil

}

func GenerateEncryptionKey() ([]byte, error) {
	// add file path to check for existing enc_key existence

	// if file empty/doesn't exist, generate a new encryption key
	encryptionKey := make([]byte, 32)

	if _, err := rand.Read(encryptionKey); err != nil {
		return nil, err
	}

	return encryptionKey, nil
}

// Seal encryption for an added layer of protection
func SealEncryptionKey(hashedPassword []byte, encryptionKey []byte) (*[24]byte, []byte, error) {
	// make slice to store rand's generated output
	generateNonce := make([]byte, 24)

	// generate nonce to use with secretbox.Seal
	if _, err := rand.Read(generateNonce); err != nil {
		return nil, nil, err
	}

	nonce := (*[24]byte)(generateNonce)

	// use secretbox to seal the encryption key
	sealedEncKey := secretbox.Seal(make([]byte, 32), encryptionKey, nonce, (*[32]byte)(hashedPassword))

	// returning nonce as well since we'll be writing the nonce, sealed encryption key
	// and the salt generated with master password to disk for subsequent usage
	return nonce, sealedEncKey, nil
}

package auth

import (
	"crypto/rand"
	"fmt"
	"os"
	"strings"

	"github.com/good-times-ahead/password-manager-go/app"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/nacl/secretbox"
)

// This file contains functions that will be executed if it is the users' run of the app
var encInfoPath = "./encrypted_data"

// Ask the user to set the master password
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
	usrInput := ""

	for true {

		prompt := "Enter desired Master Password(should contain a combination of atleast 1 lowercase, 1 uppercase letter and a number; minimum length: 8 characters): "

		usrInput = app.GetInput(prompt)

		checkInput := CheckPasswordStrength(usrInput)

		if !checkInput {
			fmt.Println("Doesn't match required parameters! Please try again.")
			fmt.Println("")
		}
		break
	}

	return usrInput

}

// CheckPasswordStrength checks the strength of the user-entered input for master password
func CheckPasswordStrength(usrInput string) bool {

	lowercase := "abcdefghijklmnopqrstuvwxyz"
	uppercase := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	nums := "0123456789"

	switch true {

	case !strings.ContainsAny(lowercase, usrInput):
		break
	case !strings.ContainsAny(uppercase, usrInput):
		break
	case !strings.ContainsAny(nums, usrInput):
		break
	case len(usrInput) < 8:
		break
	default:
		return true

	}

	return false
}

// Argon2 is considered better than bcrypt for securing passwords
func HashMasterPassword(usrInput, pwFilePath string) ([]byte, []byte, error) {
	// generate 32 bytes salt
	salt := make([]byte, 32)

	// generate 32 byte long random salt
	// rand.Read writes the output to the given slice
	if _, err := rand.Read(salt); err != nil {
		return nil, nil, err
	}

	// generate hashed password using argon2
	hashedMasterPassword := argon2.IDKey([]byte(usrInput), salt, 1, 64*1024, 4, 32)

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

// Seal encryption key for an added layer of protection
func SealEncryptionKey(hashedPassword []byte, encryptionKey []byte) ([]byte, error) {

	// make slice to store rand's generated output
	// generateNonce := make([]byte, 24)
	var nonce [24]byte

	// generate nonce to use with secretbox.Seal
	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, err
	}

	// use secretbox to seal the encryption key
	sealedEncKey := secretbox.Seal(nonce[:], encryptionKey, &nonce, (*[32]byte)(hashedPassword))

	// returning nonce as well since we'll be writing the nonce, sealed encryption key
	// and the salt generated with master password to disk for subsequent usage
	return sealedEncKey, nil

}

// Save the salt, sealed encryption key and nonce to file
func SaveEncryptionData(values [][]byte) error {

	file, err := os.OpenFile(encInfoPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0444)

	if err != nil {
		return err
	}

	for _, value := range values {

		if _, err := file.Write(append(value, '\n')); err != nil {
			return err
		}

	}

	return nil

}

func FirstRun(pwFilePath string) error {

	// generate encryption key at the very start
	encKey, err := GenerateEncryptionKey()

	if err != nil {
		return err
	}

	// if password doesn't exist yet
	usrInput := GetMasterPassword()

	salt, hashedPassword, err := HashMasterPassword(usrInput, pwFilePath)

	if err != nil {
		return err
	}

	savePasswordErr := SaveMasterPassword(pwFilePath, hashedPassword)

	if savePasswordErr != nil {
		return savePasswordErr
	}

	// after master password has been generated properly, we will seal our encryption key
	sealedEncKey, err := SealEncryptionKey(hashedPassword, encKey)

	// combine nonce, sealed key and salt into a slice

	values := [][]byte{salt, sealedEncKey}

	saveDataErr := SaveEncryptionData(values)

	if saveDataErr != nil {
		return err
	}

	return nil

}

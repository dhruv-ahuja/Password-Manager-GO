package auth

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/good-times-ahead/password-manager-go/store"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/nacl/secretbox"
)

func CheckEncryptedData(encInfoPath string) bool {

	// Check for file's existance, returns an error if unable to open
	checkFile, err := os.Open(encInfoPath)

	if err != nil {
		return false
	}

	defer checkFile.Close()

	readBytes, err := checkFile.Read(make([]byte, 64))

	// If not, return false since either the file has been tampered with
	if err != nil || readBytes == 0 {
		return false
	}

	return true

}

// Check whether master password file exists already
func CheckMasterPassword(pwFilePath string) bool {

	// Check for file's existence, returns an error if unable to open
	checkFile, err := os.Open(pwFilePath)

	if err != nil {
		return false
	}

	defer checkFile.Close()

	// Read file to confirm there are 32 bytes of the hashed master password
	readBytes, err := checkFile.Read(make([]byte, 32))

	// If not, return false since either the file has been tampered with
	if err != nil || readBytes == 0 {
		return false
	}

	return true

}

func LoadEncryptedInfo(encInfoPath string) ([][]byte, error) {

	file, err := os.Open(encInfoPath)

	if err != nil {
		return nil, fmt.Errorf("error reading encrypted data: %s", err)
	}

	// defer file.Close()

	scanner := bufio.NewScanner(file)

	count := 0

	var salt, sealedKey []byte

	// order: 1. salt, 2. sealedkey
	for scanner.Scan() {

		count += 1

		switch count {

		case 1:
			salt = []byte(scanner.Text())
		case 2:
			sealedKey = []byte(scanner.Text())
		default:

		}

	}

	values := [][]byte{salt, sealedKey}

	// returning file.Close() since the function will return an error. Deferring it means ignoring any errors that might occur, which can lead to data loss while the program continues functioning under the assumption that everything went well.
	return values, file.Close()
}

// Take users' master password input and compare it to the stored hash, allowing access if match
func AuthorizeUser(pwFilePath string, values [][]byte) error {

	// Load master password file
	hashedPassword, err := os.ReadFile(pwFilePath)

	if err != nil {
		return fmt.Errorf("error reading master password: %s", err)
	}

	// infinite loop till user enters the correct value
	for {

		// Take user input
		prompt := "\nEnter the Master Password:"
		usrPassword := store.GetPassInput(prompt)
		// convert the received slice of bytes to string
		usrInput := string(usrPassword)

		compare := argon2.IDKey([]byte(usrInput), values[0], 1, 64*1024, 4, 32)

		// if the stored hash matches the produced/current hash, allow the user to go through
		if bytes.Equal(compare, hashedPassword) {
			break
		}

		// adding new lines to keep the interface clean and readable
		fmt.Printf("\nThe passwords do not match! Try again.\n\n")

	}

	return nil

}

func UnsealEncryptionKey(pwFilePath string, values [][]byte) ([]byte, error) {

	// declare needed encryption variables from the slice
	var sealedKey []byte
	var nonce [24]byte

	for index, value := range values {

		if index == 1 {
			sealedKey = value
		}

	}

	// the nonce is stored in the first 24 bytes
	copy(nonce[:], sealedKey[:24])

	// read master password hash to use as the key for secretbox.Open
	hashedPassword, err := os.ReadFile(pwFilePath)

	if err != nil {
		return nil, fmt.Errorf("error when trying to read master password: %s", err)
	}

	// the main data is stored from the 25th byte onwards
	encKey, ok := secretbox.Open(nil, sealedKey[24:], &nonce, (*[32]byte)(hashedPassword))

	if !ok {
		return nil, errors.New("error while unsealing encryption key")
	}

	return encKey, nil
}

func Run(encInfoPath, pwFilePath string) ([]byte, error) {

	// load encrypted data to use when dealing with credentials later
	encData, err := LoadEncryptedInfo(encInfoPath)
	if err != nil {
		return nil, err
	}

	err = AuthorizeUser(pwFilePath, encData)
	if err != nil {
		return nil, err
	}

	// now, to unseal the encryption key
	encryptionKey, err := UnsealEncryptionKey(pwFilePath, encData)

	if err != nil {
		return nil, err
	}

	return encryptionKey, nil

}

package auth

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/good-times-ahead/password-manager-go/app"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/nacl/secretbox"
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

func LoadEncryptedInfo(encInfoPath string) ([][]byte, error) {

	file, err := os.Open(encInfoPath)

	if err != nil {
		return nil, errors.New("error reading from encrypted data file, please check")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	count := 0

	var salt, sealedKey, nonce []byte

	// order: 1. salt, 2. sealedkey, 3. nonce
	for scanner.Scan() {

		count += 1

		switch count {

		case 1:
			salt = []byte(scanner.Text())
		case 2:
			sealedKey = []byte(scanner.Text())
		case 3:
			nonce = []byte(scanner.Text())
		default:
			break

		}

	}

	values := [][]byte{salt, sealedKey, nonce}

	return values, nil
}

// Take users' master password input and compare it to the stored hash, allowing access if match
func AuthorizeUser(pwFilePath string, values [][]byte) error {

	// Load master password file
	hashedPassword, err := os.ReadFile(pwFilePath)

	if err != nil {
		return errors.New("error when trying to read master password file, please check")
	}

	// infinite loop till user enters the correct value
	for true {

		// Take user input
		prompt := "Enter the Master Password: "
		usrInput := app.GetInput(prompt)

		compare := argon2.IDKey([]byte(usrInput), values[0], 1, 64*1024, 4, 32)

		// if the stored hash matches the produced/current hash, allow the user to go through
		if bytes.Equal(compare, hashedPassword) {
			return nil
		}

		fmt.Println("The passwords do not match! Try again.")
		fmt.Println()

	}

	return nil

}

func UnsealEncryptionKey(pwFilePath string, values [][]byte) ([]byte, error) {

	// declare needed encryption variables from the slice
	var sealedKey []byte
	var nonce [24]byte

	for index, value := range values {

		switch index {

		case 1:
			sealedKey = value

		default:
			break

		}

	}

	// the nonce is stored in the first 24 bytes
	copy(nonce[:], sealedKey[:24])

	// read master password hash to use as the key for secretbox.Open
	hashedPassword, err := os.ReadFile(pwFilePath)

	if err != nil {
		return nil, errors.New("error when trying to read master password file, please check")
	}

	encKey, ok := secretbox.Open(nil, sealedKey[24:], &nonce, (*[32]byte)(hashedPassword))

	if !ok {
		return nil, errors.New("error while unsealing encryption key")
	}

	return encKey, nil
}

func Run(pwFilePath string) error {

	values, err := LoadEncryptedInfo(encInfoPath)
	if err != nil {
		return err
	}

	err = AuthorizeUser(pwFilePath, values)
	if err != nil {
		return err
	}

	UnsealEncryptionKey(pwFilePath, values)

	return nil

}

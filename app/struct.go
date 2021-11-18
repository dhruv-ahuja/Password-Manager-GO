package app

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/good-times-ahead/password-manager-go/database"
)

// Struct to store all user information
type credentials struct {
	ID                                 int
	username, email, website, password string
}

// Encrypt the password to allow storing it into the database safely
func (c credentials) EncryptPassword(encryptionKey []byte) (string, error) {

	password := []byte(c.password)

	// retrieve encryption key from .env file
	// key := []byte(os.Getenv("ENC_KEY"))

	// if len(key) != 32 {
	// 	return "", errors.New("'ENC_KEY' environment variable not defined properly")
	// }

	//generate a new cipher using our 32 byte long key
	generatedCipher, err := aes.NewCipher(encryptionKey)

	if err != nil {
		return "", errors.New("error generating cipher")
	}

	// GCM is a mode of operation used on block ciphers
	gcm, err := cipher.NewGCM(generatedCipher)

	if err != nil {
		return "", errors.New("error generating GCM")
	}

	nonce := make([]byte, gcm.NonceSize())

	//populates our byte array with a cryptographically secure sequence
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errors.New("error generating secured-sequence")
	}

	encryptedPassword := gcm.Seal(nonce, nonce, password, nil)

	// convert the slice of encrypted bytes into a base64 encrypted string to store in the database
	b64Password := base64.StdEncoding.EncodeToString(encryptedPassword)

	return b64Password, nil
}

// Reads all struct fields and inserts them into the database
func (c credentials) InsertIntoDB(encryptedPassword string) error {

	query := "INSERT INTO info (website, email, username, encrypted_pw) VALUES ($1, $2, $3, $4) RETURNING *"

	_, err := database.DB.Exec(query, c.website, c.email, c.username, encryptedPassword)

	if err != nil {
		return errors.New("unable to save your credentials to the database")
	}

	fmt.Println("Saved your credentials to the database!")

	return nil

}

func (c credentials) DecryptPassword(base64Password string, encryptionKey []byte) (string, error) {

	// decrypting the base64 password string to retrieve our AES-encrypted password
	encryptedPassword, err := base64.StdEncoding.DecodeString(base64Password)

	if err != nil {
		return "", errors.New("error decoding base64 encrypted password string")
	}

	// generate a new cipher using our 32 byte long key
	generatedCipher, err := aes.NewCipher(encryptionKey)

	if err != nil {
		fmt.Println(err)
		return "", errors.New("error generating cipher")
	}

	gcm, err := cipher.NewGCM(generatedCipher)

	if err != nil {
		return "", errors.New("error generating GCM")
	}

	nonce, ciphertext := encryptedPassword[:gcm.NonceSize()], encryptedPassword[gcm.NonceSize():]

	password, err := gcm.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		fmt.Println(err)
		return "", errors.New("error attempting to decrypt AES-encrypted password")
	}

	return string(password), nil

}

// Update credentials using ID number
func (c credentials) UpdateCredentials(modifyPassword bool) error {
	// Since the password is the key component here, we specifically set a flag for it
	// Update password as well if the bool is true otherwise only update username and email
	if modifyPassword {
		query := "UPDATE info SET email = $1, username = $2, encrypted_pw = $3 WHERE id=$4"

		_, err := database.DB.Exec(query, c.email, c.username, c.password, c.ID)

		if err != nil {
			return errors.New("error executing query")
		}
		// Skip updating the password otherwise
	} else {
		query := "UPDATE info SET email = $1, username = $2 WHERE id=$3"

		_, err := database.DB.Exec(query, c.email, c.username, c.ID)

		if err != nil {
			return errors.New("error executing query")
		}
	}

	return nil

}

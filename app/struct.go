package app

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/good-times-ahead/password-manager-go/database"
	"github.com/joho/godotenv"
)

// Struct to store all user information
type credentials struct {
	ID                                 int
	username, email, website, password string
}

// Encrypt the password to allow storing it into the database safely
func (c credentials) EncryptPassword() (string, error) {

	if err := godotenv.Load(); err != nil {
		return "", err
	}

	password := []byte(c.password)
	key := []byte(os.Getenv("ENC_KEY"))

	//generate a new cipher using our 32 byte long key
	generatedCipher, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// GCM is a mode of operation for block ciphers
	gcm, err := cipher.NewGCM(generatedCipher)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())

	//populates our byte array with a cryptographically secure sequence
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	encryptedPassword := gcm.Seal(nonce, nonce, password, nil)

	b64Password := base64.StdEncoding.EncodeToString(encryptedPassword)

	return b64Password, nil
}

func (c credentials) InsertIntoDB(encryptedPassword string) error {

	query := "INSERT INTO info (website, email, username, encrypted_pw) VALUES ($1, $2, $3, $4) RETURNING *"

	row := database.DB.QueryRow(query, c.website, c.email, c.username, encryptedPassword)

	var usrInfo credentials

	if err := row.Scan(&usrInfo.ID, &usrInfo.website, &usrInfo.email, &usrInfo.username, &usrInfo.password); err != nil {

		if err == sql.ErrNoRows {

			return err
		}

	}

	response := fmt.Sprintf("ID: %d, Website: %s, Email: %s, Username: %s, Encrypted Password: %s", usrInfo.ID, usrInfo.website, usrInfo.email, usrInfo.username, usrInfo.password)

	fmt.Println(response)

	fmt.Println("Saved your credentials to the database!")
	return nil

}

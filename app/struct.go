package app

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
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

	fmt.Println("encrypted password is: %s", encryptedPassword)

	return "", err
}

func (c credentials) InsertIntoDB(encryptedPassword string) error {

	// Preparing 1st half of the SQL query
	query := fmt.Sprintf("INSERT INTO %s (website, email, username, password_hash)", database.Table)

	_, err := database.DB.Exec(query+"VALUES ($1, $2, $3, $4)", c.website, c.email, c.username, encryptedPassword)

	if err != nil {
		return err
	}

	fmt.Println("Saved your credentials to the database!")
	return nil

}

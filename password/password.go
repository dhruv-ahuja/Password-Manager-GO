package password

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/good-times-ahead/password-manager-go/credentials"
)

// Encrypt the plain-text password
func Encrypt(encryptionKey []byte, c credentials.Credentials) (string, error) {

	password := []byte(c.Password)

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

// Decrypt the base64 encoded password
func Decrypt(base64Password string, encryptionKey []byte, c credentials.Credentials) (string, error) {

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

package password

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// Encrypt the plain-text password
func Encrypt(encryptionKey []byte, plainText string) (string, error) {

	password := []byte(plainText)

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

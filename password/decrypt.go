package password

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
)

// Decrypt the base64 encoded password
func Decrypt(base64Password string, encryptionKey []byte) (string, error) {

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

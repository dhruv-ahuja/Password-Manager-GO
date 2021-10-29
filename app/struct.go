package app

import (
	"fmt"

	"github.com/good-times-ahead/password-manager-go/database"
	"golang.org/x/crypto/bcrypt"
)

// Struct to store all user information
type credentials struct {
	username, email, website, password string
}

// Salt and hash the password to allow storing it in the database safely
func (c credentials) HashPassword() (string, error) {

	password := c.password
	// convert password(string) to a slice of encryptedPassword for hashing & salting.
	// Cost(from what I understand) is the number of protective layers to add to the password
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(encryptedPassword), nil
	// decryptedPass := bcrypt.CompareHashAndPassword(encryptedPassword, []byte(password))

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

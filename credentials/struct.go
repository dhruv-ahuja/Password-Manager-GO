package credentials

import (
	"errors"
	"fmt"

	"github.com/good-times-ahead/password-manager-go/database"
)

// Struct to store all user information
type Credentials struct {
	ID                                 int
	Username, Email, Website, Password string
}

// Reads all struct fields and inserts them into the database
func (c Credentials) InsertIntoDB(encryptedPassword string) error {

	query := "INSERT INTO info (website, email, username, encrypted_pw) VALUES ($1, $2, $3, $4) RETURNING *"

	_, err := database.DB.Exec(query, c.Website, c.Email, c.Username, encryptedPassword)

	if err != nil {
		return errors.New("unable to save your credentials to the database")
	}

	fmt.Println("Saved your credentials to the database!")

	return nil

}

// Update credentials using ID number
func (c Credentials) UpdateCredentials(modifyPassword bool) error {
	// Since the password is the key component here, we specifically set a flag for it
	// Update password as well if the bool is true otherwise only update username and email
	if modifyPassword {
		query := "UPDATE info SET email = $1, username = $2, encrypted_pw = $3 WHERE id=$4"

		_, err := database.DB.Exec(query, c.Email, c.Username, c.Password, c.ID)

		if err != nil {
			return errors.New("error executing query")
		}
		// Skip updating the password otherwise
	} else {
		query := "UPDATE info SET email = $1, username = $2 WHERE id=$3"

		_, err := database.DB.Exec(query, c.Email, c.Username, c.ID)

		if err != nil {
			return errors.New("error executing query")
		}
	}

	return nil

}

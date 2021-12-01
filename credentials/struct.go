package credentials

import (
	"errors"
	"fmt"

	"github.com/good-times-ahead/password-manager-go/database"
)

type CredentialFuncs interface {
	InsertIntoDB(string) error

	UpdateCredentials(bool) error
}

// Struct to store all user information
type Credentials struct {
	ID            int
	Key, Password string
	Repo          database.Repo
}

// Reads all struct fields and inserts them into the database
func (c Credentials) InsertIntoDB(encryptedPassword string) error {
	//TODO: return the result of the executed query
	query := "INSERT INTO info (key, encrypted_pw) VALUES ($1, $2) RETURNING *"

	_, err := c.Repo.DB.Exec(query, c.Key, encryptedPassword)

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
		query := "UPDATE info SET key = $1, encrypted_pw = $2  WHERE id= $3"

		_, err := c.Repo.DB.Exec(query, c.Key, c.Password, c.ID)

		if err != nil {
			return errors.New("error executing query")
		}
		// Skip updating the password otherwise
	} else {
		query := "UPDATE info SET key = $1 WHERE id = $2"

		_, err := c.Repo.DB.Exec(query, c.Key, c.ID)

		if err != nil {
			return errors.New("error executing query")
		}
	}

	return nil

}

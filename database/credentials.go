package database

import (
	"errors"
	"fmt"
)

// This file contains functions that interact with the DB to
// either save or update user credentials

type CredentialFuncs interface {
	InsertIntoDB(string, map[string]string) error

	UpdateCredentials(bool, map[string]string) error
}

// Reads all struct fields and inserts them into the database
func (conn *Repo) InsertIntoDB(encryptedPassword string, credentials map[string]string) error {

	//TODO: add credentials parameter of type map containing all necessary data
	// TODO: return the result of the executed query
	query := "INSERT INTO info (key, encrypted_pw) VALUES ($1, $2) RETURNING *"

	_, err := conn.DB.Exec(query, credentials["key"], encryptedPassword)

	if err != nil {
		return errors.New("unable to save your credentials to the database")
	}

	fmt.Println("Saved your credentials to the database!")

	return nil

}

// Update credentials using ID number
func (conn *Repo) UpdateCredentials(modifyPassword bool, credentials map[string]string) error {

	// Update password only if it has been modified
	if modifyPassword {
		query := "UPDATE info SET key = $1, encrypted_pw = $2  WHERE id= $3"

		_, err := conn.DB.Exec(query, credentials["key"], credentials["password"], credentials["id"])

		if err != nil {
			return errors.New("error executing query")
		}
		// Skip updating the password otherwise
	} else {
		query := "UPDATE info SET key = $1 WHERE id = $2"

		_, err := conn.DB.Exec(query, credentials["key"], credentials["id"])

		if err != nil {
			return errors.New("error executing query")
		}
	}

	return nil

}

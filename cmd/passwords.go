package cmd

import (
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Struct storing all user information
type credentials struct {
	username, email, website, password string
}

// Salt and hash the password to allow storing it in database safely
func (credentials) HashPassword(c credentials) (string, error) {

	password := c.password

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	return string(bytes), err

}

// Take user input
func TakeInput(db *sql.DB, table string) error {

	mainMsg := `Hello, what would you like to do today?
	1. Save a password to the DB.
	2. View a saved password.
	3. Edit a saved password.
	`
	fmt.Println(mainMsg)

	var usrInput string
	fmt.Scanf("%s", &usrInput)

	if usrInput == "1" {
		return SavePassword(db, table)
	}

	return nil
}

// todo: implement a "mother" function which will guide how each function works, just like in series-lookup proj.

// Take user information and save it to the database.
func SavePassword(db *sql.DB, table string) error {

	// defining needed prompts
	promptWebsite := "Enter the websites' name: "
	promptEmail := "Enter your accounts' email: "
	promptUsername := "Enter your username: "
	promptPassword := "Enter your password(it will be hashed before saving): "

	// making variables to store userInfo
	var usrInfo credentials

	// Scan automatically ensures that the input isn't empty
	fmt.Print(promptWebsite)
	fmt.Scan(&usrInfo.website)

	fmt.Print(promptEmail)
	fmt.Scan(&usrInfo.email)

	fmt.Print(promptUsername)
	fmt.Scan(&usrInfo.username)

	fmt.Print(promptPassword)
	fmt.Scan(&usrInfo.password)

	hashedPassword, hashingErr := usrInfo.HashPassword(usrInfo)
	if hashingErr != nil {
		return hashingErr
	}

	//if no errors occured, time to write the credentials to the database
	if saveToDB := InsertIntoDB(usrInfo, db, table, hashedPassword); saveToDB != nil {
		return errors.New("problem saving to DB")
	}

	return nil
}

// Write credentials to database
func InsertIntoDB(c credentials, db *sql.DB, table, hashedPassword string) error {
	// Execute a query feeding credentials to the databases
	// First prepare that part of the statement containing a variable value to be inserted since we can't do it dynamically when running the query itself
	query := fmt.Sprintf("INSERT INTO %s (website, email, username, password_hash)", table)

	_, err := db.Exec(query+"VALUES ($1, $2, $3, $4)", c.website, c.email, c.username, hashedPassword)

	if err != nil {
		return err
	}

	return nil
}

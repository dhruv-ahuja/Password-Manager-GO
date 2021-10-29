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

// Take user input to dictate what function gets executed
func TakeInput() error {

	mainMsg := `Hello, what would you like to do today?
1. Save a password to the DB.
2. View a saved password.
3. Edit a saved password.`
	fmt.Println(mainMsg)

	var usrInput string
	fmt.Scanf("%s", &usrInput)

	if usrInput == "1" {
		// returning function directly since it's supposed to return an error anyway
		return SaveCredentials()
	}

	return nil
}

// Save credentials to the database
func SaveCredentials() error {

	// define needed prompts
	promptWebsite := "Enter the websites' name: "
	promptEmail := "Enter your mail ID: "
	promptUsername := "Enter your username: "
	promptPassword := "Enter your password(it will be encrypted before saving): "

	// initialize the variable to save the credentials to
	var usrInfo credentials

	// Scan automatically ensures that the input isn't empty
	fmt.Println(promptWebsite)
	fmt.Scan(&usrInfo.website)

	fmt.Println(promptEmail)
	fmt.Scan(&usrInfo.email)

	fmt.Println(promptUsername)
	fmt.Scan(&usrInfo.username)

	fmt.Println(promptPassword)
	fmt.Scan(&usrInfo.password)

	encryptedPassword, err := usrInfo.HashPassword()
	if err != nil {
		return err
	}

	// save the credentials to the database
	if saveToDB := usrInfo.InsertIntoDB(encryptedPassword); saveToDB != nil {
		return saveToDB
	}

	return nil

}

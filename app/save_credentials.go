package app

import "fmt"

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

	encryptedPassword, err := usrInfo.EncryptPassword()
	if err != nil {
		return err
	}

	// save the credentials to the database
	if saveToDB := usrInfo.InsertIntoDB(encryptedPassword); saveToDB != nil {
		return saveToDB
	}

	return nil

}

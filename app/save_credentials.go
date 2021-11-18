package app

// Save credentials to the database
func SaveCredentials(encryptionKey []byte) error {

	// define needed prompts
	promptWebsite := "Enter the websites' name: "
	promptEmail := "Enter your mail ID: "
	promptUsername := "Enter your username: "
	promptPassword := "Enter your password(it will be encrypted before saving): "

	// initialize the variable to save the credentials to
	var usrInfo credentials

	// Write user input to respective structure fields
	usrInfo.website = GetInput(promptWebsite)
	usrInfo.email = GetInput(promptEmail)
	usrInfo.username = GetInput(promptUsername)
	usrInfo.password = GetInput(promptPassword)

	// encrypt the plain text password
	encryptedPassword, err := usrInfo.EncryptPassword(encryptionKey)

	if err != nil {
		return err
	}

	// save the credentials to the database
	if saveToDB := usrInfo.InsertIntoDB(encryptedPassword); saveToDB != nil {
		return saveToDB
	}

	return nil

}

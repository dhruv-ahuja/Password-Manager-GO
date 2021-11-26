package app

// View credentials for the specified website
func ViewSavedCredentials(key string, encryptionKey []byte) error {

	// Get all accounts associated with the website
	query := "SELECT * FROM info WHERE key = $1 ORDER BY id ASC"

	_, err := RetrieveCredentials(query, key, encryptionKey)

	if err != nil {
		return err
	}

	return nil

}

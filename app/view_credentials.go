package app

// View credentials for the specified website
func ViewSavedCredentials(website string) error {
	// Get all accounts associated with the website
	query := "SELECT * FROM info WHERE website=$1 ORDER BY id ASC"

	_, err := RetrieveCredentials(query, website)

	if err != nil {
		return err
	}

	return nil

}

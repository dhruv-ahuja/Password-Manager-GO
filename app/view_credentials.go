package app

// View credentials for the specified website
func (DBConn *DBConn) ViewCredentials(key string, encryptionKey []byte) error {

	// Get all accounts associated with the website
	query := "SELECT * FROM info WHERE key ILIKE $1 ORDER BY id ASC;"

	_, err := DBConn.RetrieveCredentials(query, key, encryptionKey)

	if err != nil {
		return err
	}

	return nil

}

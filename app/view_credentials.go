package app

// View credentials for the specified website
// func ViewCredentials(website string) error {
// 	// Get all accounts associated with the website
// 	query := fmt.Sprintf("SELECT * FROM %s", database.Table)

// 	rows, err := database.DB.Query(query+"WHERE website=$1", website)
// 	if err != nil {
// 		return err
// 	}

// 	defer rows.Close()

// 	// making a slice to store all rows' data
// 	var accounts []credentials

// 	for rows.Next() {
// 		var usrInfo credentials
// 		// We need to decrypt the password before we can print it out to the user
// 		var encryptedPassword string
// 		// Scan takes pointers to variables to write data
// 		if err := rows.Scan(&usrInfo.ID, &usrInfo.website, &usrInfo.email, &usrInfo.username, &encryptedPassword); err != nil {
// 			return err
// 		}

// 	}

// }

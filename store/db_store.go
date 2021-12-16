package Store

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/good-times-ahead/password-manager-go/password"
)

type DBStore struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
}

// NewDBStore returns a new DBStore instance containing
// all necessary connection parameters
func NewDBStore() DBStore {

	return DBStore{
		host:     os.Getenv("HOST"),
		port:     os.Getenv("PORT"),
		user:     os.Getenv("DB_USER"),
		password: os.Getenv("DB_PASSWORD"),
		dbname:   os.Getenv("DB_NAME"),
	}
}

func (db DBStore) NewConnection() (*sql.DB, error) {
	// Prepare Postgres connection parameters
	psqlInfo := fmt.Sprint("host=", db.host, " port=", db.port,
		" user=", db.user, " password=", db.password,
		" dbname=", db.dbname, " sslmode=disable")

	conn, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		fmt.Println("error connecting to database!")
		return nil, err
	}

	// Ping to confirm whether connection works
	if err := conn.Ping(); err != nil {
		fmt.Println("error pinging database!")
		return nil, err
	}

	fmt.Println("Connected to the database successfully!")

	return conn, nil

}

// MakeTable creates the table to be used for our database-related operations
func (db DBStore) CreateTable(sqlFilePath string, conn *sql.DB) error {

	fmt.Println("First-time execution; creating table...")

	// Read the file content
	query, err := os.ReadFile(sqlFilePath)

	if err != nil {
		fmt.Println("unable to read setup.sql file!")
		return err
	}

	if _, err := conn.Exec(string(query)); err != nil {
		fmt.Println("unable to create table!")
		return err
	}

	// adding new lines to keep interface clean & readable
	fmt.Printf("Everything done. You're good to go.\n\n")

	return nil

}

// Interface functions

// SaveCreds saves the user-entered credentials to the database
func (db DBStore) SaveCreds(encryptionKey []byte, conn *sql.DB) error {
	// define needed prompts
	promptKey := "Enter key: "
	promptPassword := "Enter your password(it will be encrypted before saving): "

	// initialize the variable to save the creds to
	usrInfo := make(map[string]string, 3)

	// write user input to the map/dict
	usrInfo["key"] = GetInput(promptKey)
	usrInfo["password"] = string(GetPassInput(promptPassword))

	// encrypt the plaintext password
	encryptedPassword, err := password.Encrypt(encryptionKey, usrInfo["password"])

	if err != nil {
		return err
	}

	// save the credentials to the database
	if err := db.InsertIntoDB(encryptedPassword, usrInfo, conn); err != nil {
		return err
	}

	return nil

}

// RetrieveCreds retrieves userdata from the database given respective query
func (db DBStore) RetrieveCreds(query, key string, encryptionKey []byte, conn *sql.DB) ([]map[string]string, error) {

	// implementing ILIKE search using the 2 "%" signs
	rows, err := conn.Query(query, "%"+key+"%")

	if err != nil {
		return nil, fmt.Errorf("error executing query: %s", err)
	}

	// prepare a slice of maps to store retrieved credentials
	var credList []map[string]string

	for rows.Next() {
		// each fetched row gets its own "usrInfo" map of length 3
		usrInfo := make(map[string]string, 3)
		var id, key, base64Password string

		// write scanned values using the variable's pointers
		err := rows.Scan(&id, &key, &base64Password)

		if err != nil {
			return nil, fmt.Errorf("error reading query data: %s", err)

		}

		// now, to decrypt the b64 string
		password, err := password.Decrypt(base64Password, encryptionKey)

		if err != nil {
			return nil, err
		}

		// finally, we populate the map we created at the beginning of the loop
		usrInfo["id"], usrInfo["key"], usrInfo["password"] = id, key, password

		// Append the map to the slice of maps
		credList = append(credList, usrInfo)
	}

	// no such key in db if slice is empty
	if len(credList) == 0 {
		fmt.Println("Nothing found for that search entry!")
	} else {
		printEntries(credList)
	}

	return credList, nil
}

func (db DBStore) ViewCreds(key string, encryptionKey []byte, conn *sql.DB) error {
	// get all accounts associated with the website
	query := "SELECT * FROM info WHERE key ILIKE $1 ORDER BY id ASC;"

	_, err := db.RetrieveCreds(query, key, encryptionKey, conn)

	if err != nil {
		return err
	}

	return nil

}

func (db DBStore) EditCreds(key string, encryptionKey []byte, conn *sql.DB) error {

	query := "SELECT * FROM info WHERE key ILIKE $1 ORDER BY id ASC;"

	// retrieve the saved account list first
	credList, err := db.RetrieveCreds(query, key, encryptionKey, conn)

	if err != nil {
		return err
	}

	// if no accounts found,
	if len(credList) == 0 {
		return nil
	}

	selectID := false
	usrInput := ""
	selection := make(map[string]string, 3)

	for !selectID {

		// Get users' input to find the entry they want to modify
		msg := "Enter the ID No. of the entry you want to modify: "

		usrInput = GetInput(msg)

		for _, entry := range credList {

			if entry["id"] == usrInput {
				selectID = true
				selection = entry
				break
			}
		}

		if !selectID {
			fmt.Println("Entered ID outside range!")
		}

	}

	// now, we have the users' choice of entry, allow them to edit whatever field they want
	fmt.Println("Your current 'key' entry is: ", selection["key"])

	fmt.Println("Enter new key (leave field blank if no changes): ")

	// using bufio since GetInput doesn't allow empty input
	reader := bufio.NewReader(os.Stdin)

	newKey, err := reader.ReadString('\n')

	if err != nil {
		return err
	}

	// trim away whitespace left over by reader
	newKey = strings.TrimSpace(newKey)

	// update the map if newKey has been modified else let it be
	if newKey != "" {
		selection["key"] = newKey
	}

	fmt.Println("\nYour current password is: ", selection["password"])

	newPassPrompt := "Enter new password: "

	newPassword := string(GetPassInput(newPassPrompt))

	// prepare to encrypt password
	b64Password, err := password.Encrypt(encryptionKey, newPassword)

	if err != nil {
		return err
	}

	// now, update the selection dict/map and send it to the database
	selection["password"] = b64Password

	db.UpdateCreds(selection, conn)

	fmt.Println("Updated your credentials successfully!")

	return nil
}

// Now adding functions that shall help us insert new and update exisitng data in the table

func (db DBStore) InsertIntoDB(encryptedPassword string, creds map[string]string, conn *sql.DB) error {

	query := "INSERT INTO info (key, encrypted_pw) VALUES ($1, $2) RETURNING (id, key)"

	var id int
	var key string

	output := conn.QueryRow(query, creds["key"], encryptedPassword)

	switch err := output.Scan(&id, &key); err {
	case sql.ErrNoRows:
		return fmt.Errorf("No rows returned by query: %s", err)

	case nil:
		fmt.Println("Saved your credentials to the database!")

		fmt.Printf("\nID: %d, Key: %s", id, key)

	default:
		return fmt.Errorf("unable to insert into DB: %s", err)
	}

	fmt.Println()

	return nil

}

func (DB DBStore) UpdateCreds(creds map[string]string, conn *sql.DB) error {

	query := "UPDATE info SET key = $1, encrypted_pw = $2  WHERE id= $3"

	_, err := conn.Exec(query, creds["key"], creds["password"], creds["id"])

	if err != nil {
		return fmt.Errorf("error updating credentials: %s", err)
	}

	return nil
}

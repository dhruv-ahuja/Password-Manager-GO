package main

import (
	"fmt"
	"log"

	"github.com/good-times-ahead/password-manager-go/app"
	"github.com/good-times-ahead/password-manager-go/auth"
	"github.com/good-times-ahead/password-manager-go/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	// Path to hashed master password file & encrypted data
	pwFilePath := "./master_pw"
	encInfoPath := "./encrypted_data"
	sqlFilePath := "./database/setup.sql"

	if err := godotenv.Load(); err != nil {
		log.Fatal(
			fmt.Errorf("error reading from .env file: %s", err),
		)
	}

	dbConfig := database.NewConfig()

	dbConn, err := database.NewConnection(dbConfig)

	defer dbConn.Close()

	if err != nil {
		log.Fatal(err)
	}

	repo := database.NewDBRepo(dbConn)

	program := app.NewProgram(repo)

	checkEncData := auth.CheckEncryptedData(encInfoPath)

	if !checkEncData {

		if err := program.Repo.MakeTable(sqlFilePath); err != nil {
			log.Fatal(err)
		}

		if err := auth.FirstRun(encInfoPath, pwFilePath); err != nil {
			log.Fatal(err)
		}
	}

	encryptionKey, err := auth.Run(encInfoPath, pwFilePath)

	if err != nil {
		log.Fatal(err)
	}

	// Finally, start the app
	if err := startApp(program, encryptionKey); err != nil {
		log.Fatal(err)
	}

}

func startApp(p *app.Program, encryptionKey []byte) error {

	appPersist := true

	for appPersist {
		mainMsg := `Hello, what would you like to do?
1. Save a password to the DB
2. View a saved password
3. Edit a saved password
4. Delete a saved password
0: Exit the application: `

		usrInput := app.GetInput(mainMsg)

		if err := p.TakeInput(usrInput, encryptionKey); err != nil {
			return err
		}
		// adding new lines to keep the interface clean and readable
		fmt.Println()

	}

	return nil

}

package main

import (
	"log"

	"github.com/good-times-ahead/password-manager-go/app"
	"github.com/good-times-ahead/password-manager-go/auth"
	"github.com/good-times-ahead/password-manager-go/database"
	_ "github.com/lib/pq"
)

func main() {
	err := database.ConnecttoDB()

	if err != nil {
		log.Fatal(err)
	}

	err = database.TableExists()

	if err != nil {
		log.Fatal(err)
	}

	checkPassword := auth.CheckMasterPassword()

	if checkPassword != nil {
		log.Fatal(checkPassword)
	}

	run := app.TakeInput()

	if run != nil {
		log.Fatal(run)
	}

}

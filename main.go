package main

import (
	"log"

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

}

package Store

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// PrintEntries outputs data fetched from queries
func printEntries(acctList []map[string]string) {

	for _, usrInfo := range acctList {

		response := fmt.Sprintf("ID No: %s, Key: %s, Password: %s", usrInfo["id"], usrInfo["key"], usrInfo["password"])

		fmt.Println(response)
		fmt.Println() // empty line to clean up the output
	}
}

// GetInput receives user input in a streamlined fashion.
func GetInput(argument string) string {

	// Emulate a while loop to receive user input and ensure its' validity
	reader := bufio.NewReader(os.Stdin)

	isEmpty := true

	for isEmpty {

		fmt.Print(argument)

		usrInput, err := reader.ReadString('\n')

		// We want to keep re-iterating over the loop so we can leave the error as is(I think)
		if err != nil {
			fmt.Println("Invalid input or method!")
		}

		switch len(usrInput) {

		case 1:
			fmt.Printf("Empty input!\n\n")
		default:
			fmt.Println()
			return strings.TrimSpace(usrInput)

		}

	}

	return ""

}

// GetPassInput receives user's input securely,
// it makes use of the syscall library to achieve this
func GetPassInput(argument string) []byte {

	// Emulate a while loop to receive user input and ensure its' validity
	isEmpty := true

	for isEmpty {

		fmt.Print(argument)

		usrInput, err := term.ReadPassword(int(syscall.Stdin))

		if err != nil {
			fmt.Println("Invalid input or method!")
		}

		switch len(usrInput) {

		case 0:
			fmt.Printf("Empty input!\n\n")
		default:
			return usrInput

		}

	}

	return nil

}

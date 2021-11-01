package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Function to get user input in a streamlined fashion.
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
			continue
		}

		if len(usrInput) != 0 {
			return strings.TrimSpace(usrInput)
		} else {
			fmt.Println("Empty input!")
		}

	}

	return ""

}

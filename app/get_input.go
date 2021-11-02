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
		//fmt.Println(usrInput, len(usrInput))

		// We want to keep re-iterating over the loop so we can leave the error as is(I think)
		if err != nil {
			fmt.Println("Invalid input or method!")
		}

		switch len(usrInput) {

		case 1:
			fmt.Println("Empty input!")

		default:
			return strings.TrimSpace(usrInput)
		}

		// if usrInput != " " {
		// 	return strings.TrimSpace(usrInput)
		// } else {
		// 	fmt.Println("Empty input!")
		// }

	}

	return ""

}

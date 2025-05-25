package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}
		// Remove trailing newline character
		command = command[:len(command)-1]

		evaluateCommandLine(command)
	}
}

func evaluateCommandLine(commandLine string) {
	commandLineParts := strings.Split(commandLine, " ")
	command := commandLineParts[0]

	switch command {
	case "exit":
		errorCode, parsingError := strconv.ParseInt(commandLineParts[1], 10, 32)
		if parsingError != nil {
			os.Exit(0)
		}
		os.Exit(int(errorCode))
	default:
		fmt.Println(command + ": command not found")
	}
}

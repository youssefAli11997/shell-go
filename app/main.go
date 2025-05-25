package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"slices"
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
	commandLineParts := whiteSpaceRegex.Split(commandLine, -1)
	command := commandLineParts[0]

	switch command {
	case "exit":
		if len(commandLineParts) == 1 {
			// If no argument is provided, just exit
			os.Exit(0)
		}
		errorCode, parsingError := strconv.ParseInt(commandLineParts[1], 10, 32)
		if parsingError != nil {
			os.Exit(0)
		}
		os.Exit(int(errorCode))
	case "echo":
		fmt.Println(strings.Join(commandLineParts[1:], " "))
	case "type":
		if len(commandLineParts) < 2 {
			fmt.Println()
		}
		commandsToType := commandLineParts[1:]
		for _, cmd := range commandsToType {
			if slices.Contains(builtinCommands, cmd) {
				fmt.Printf("%s is a shell builtin\n", cmd)
			} else {
				fmt.Printf("%s: not found\n", cmd)
			}
		}
	default:
		fmt.Println(command + ": command not found")
	}
}

var builtinCommands = []string{
	"exit",
	"echo",
	"type",
}

var whiteSpaceRegex = regexp.MustCompile(`\s+`)

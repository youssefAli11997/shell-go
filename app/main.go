package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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
		evaluateExit(commandLineParts)
	case "echo":
		evaluateEcho(commandLineParts)
	case "type":
		evaluateType(commandLineParts)
	case "pwd":
		evaluatePwd()
	default:
		evaluateExternalCommand(commandLineParts)
	}
}

func evaluateExit(commandLineParts []string) {
	if len(commandLineParts) == 1 {
		// If no argument is provided, just exit
		os.Exit(0)
	}
	errorCode, parsingError := strconv.ParseInt(commandLineParts[1], 10, 32)
	if parsingError != nil {
		os.Exit(0)
	}
	os.Exit(int(errorCode))
}

func evaluateEcho(commandLineParts []string) {
	fmt.Println(strings.Join(commandLineParts[1:], " "))
}

func evaluateType(commandLineParts []string) {
	if len(commandLineParts) < 2 {
		fmt.Println()
	}
	commandsToType := commandLineParts[1:]
	for _, cmd := range commandsToType {
		if slices.Contains(builtinCommands, cmd) {
			fmt.Printf("%s is a shell builtin\n", cmd)
		} else if path, err := exec.LookPath(cmd); err == nil {
			fmt.Printf("%s is %s\n", cmd, path)
		} else {
			fmt.Printf("%s: not found\n", cmd)
		}
	}
}

func evaluatePwd() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting current directory:", err)
		return
	}
	fmt.Println(currentDir)
}

func evaluateExternalCommand(commandLineParts []string) {
	command := commandLineParts[0]
	arguments := commandLineParts[1:]
	_, err := exec.LookPath(command)
	if err != nil {
		evaluateNotFoundCommand(command)
		return
	}

	cmd := exec.Command(command, arguments...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Run()
}

func evaluateNotFoundCommand(command string) {
	fmt.Printf("%s: command not found\n", command)
}

var builtinCommands = []string{
	"exit",
	"echo",
	"type",
	"pwd",
}

var whiteSpaceRegex = regexp.MustCompile(`\s+`)

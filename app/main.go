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
		// currentDir, err := os.Getwd()
		// if err != nil {
		// 	fmt.Fprintln(os.Stderr, "Error getting current directory:", err)
		// 	os.Exit(1)
		// }
		// fmt.Fprintf(os.Stdout, "(%s) $ ", currentDir)
		fmt.Fprintf(os.Stdout, "$ ")

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
	command, arguments := extractCommandAndArguments(commandLine)

	switch command {
	case "exit":
		evaluateExit(arguments)
	case "echo":
		evaluateEcho(arguments)
	case "type":
		evaluateType(arguments)
	case "pwd":
		evaluatePwd()
	case "cd":
		evaluateCD(arguments)
	default:
		evaluateExternalCommand(command, arguments)
	}
}

func evaluateExit(arguments []string) {
	if len(arguments) == 0 {
		// If no argument is provided, just exit
		os.Exit(0)
	}
	errorCode, parsingError := strconv.ParseInt(arguments[0], 10, 32)
	if parsingError != nil {
		os.Exit(0)
	}
	os.Exit(int(errorCode))
}

func evaluateEcho(arguments []string) {
	fmt.Println(strings.Join(arguments, " "))
}

func evaluateType(arguments []string) {
	if len(arguments) == 0 {
		fmt.Println()
	}
	for _, cmd := range arguments {
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

func evaluateCD(arguments []string) {
	var targetDir string
	var err error
	if len(arguments) == 0 || arguments[0] == "" {
		targetDir, err = os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, "cd: cannot find home directory")
			return
		}
	} else {
		targetDir = arguments[0]
		if strings.HasPrefix(targetDir, "~") {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Fprintln(os.Stderr, "cd: cannot find home directory")
				return
			}
			targetDir = strings.Replace(targetDir, "~", homeDir, 1)
		}
	}
	err = os.Chdir(targetDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", targetDir)
	}
}

func evaluateExternalCommand(command string, arguments []string) {
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

func extractCommandAndArguments(commandLine string) (string, []string) {
	commandLine = strings.TrimSpace(commandLine)
	var re = regexp.MustCompile(`'(\w+)'|"(\w+)"`)
	commandLine = re.ReplaceAllString(commandLine, `$1`)
	if commandLine == "" {
		return "", nil
	}
	matches := argRegex.FindAllString(commandLine, -1)
	if len(matches) == 0 {
		return "", nil
	}

	// Remove quotes from arguments
	for i, match := range matches {
		matches[i] = strings.Trim(match, `"'`)
	}

	return matches[0], matches[1:]
}

// Regex to match all arguments (quoted or unquoted) as a list of strings
var argRegex = regexp.MustCompile(`'[^']*'|"[^"]*"|\S+`)

var builtinCommands = []string{
	"exit",
	"echo",
	"type",
	"pwd",
}

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

// splitCommand splits a command into arguments.
func splitCommand(command string) []string {
	re := regexp.MustCompile(`(?:[^\s'"]+|['"][^'"]*['"])`)
	matches := re.FindAllString(command, -1)
	for i, match := range matches {
		if len(match) > 1 && (match[0] == '"' || match[0] == '\'') && match[len(match)-1] == match[0] {
			matches[i] = match[1 : len(match)-1]
		}
	}
	return matches
}

// runCommand runs a command.
func runCommand(command string) error {
	var cmd *exec.Cmd
	commandParts := splitCommand(command)
	lenCmdStringParts := len(commandParts)
	if lenCmdStringParts == 0 {
		return fmt.Errorf("invalid command: %s", command)
	} else if lenCmdStringParts == 1 {
		cmd = exec.Command(commandParts[0])
	} else {
		cmd = exec.Command(commandParts[0], commandParts[1:]...)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

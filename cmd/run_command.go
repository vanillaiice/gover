package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/mattn/go-shellwords"
)

// splitCommand splits a command into arguments.
func splitCommand(command string) ([]string, error) {
	return shellwords.Parse(command)
}

// runCommand runs a command.
func runCommand(command string) error {
	var cmd *exec.Cmd
	commandParts, err := splitCommand(command)
	if err != nil {
		return err
	}
	lenCmdStringParts := len(commandParts)

	switch lenCmdStringParts {
	case 0:
		return fmt.Errorf("invalid command: %s", command)
	case 1:
		cmd = exec.Command(commandParts[0])
	default:
		cmd = exec.Command(commandParts[0], commandParts[1:]...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

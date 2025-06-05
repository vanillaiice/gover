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

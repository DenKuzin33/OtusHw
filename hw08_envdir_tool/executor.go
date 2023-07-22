package main

import (
	"errors"
	"os"
	"os/exec"
)

var exitError *exec.ExitError

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...) //#nosec G204
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	for key, value := range env {
		if value.NeedRemove {
			os.Unsetenv(key)
		} else {
			os.Setenv(key, value.Value)
		}
	}

	result := command.Run()

	if result == nil {
		return 0
	}

	if errors.As(result, &exitError) {
		return exitError.ExitCode()
	}

	return -1
}

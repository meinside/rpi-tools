package command

import (
	"fmt"
	"os/exec"
	"strings"
)

// Run runs given command with parameters and return combined output
func Run(cmdAndParams ...string) (string, error) {
	if len(cmdAndParams) <= 0 {
		return "", fmt.Errorf("no command provided")
	}

	output, err := exec.Command(cmdAndParams[0], cmdAndParams[1:]...).CombinedOutput()
	return strings.TrimRight(string(output), "\n"), err
}

// SudoRun runs given command with parameters and return combined output (with sudo)
func SudoRun(cmdAndParams ...string) (string, error) {
	if len(cmdAndParams) <= 0 {
		return "", fmt.Errorf("no command provided")
	}

	output, err := exec.Command("sudo", cmdAndParams...).CombinedOutput()
	return strings.TrimRight(string(output), "\n"), err
}

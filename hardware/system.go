package hardware

// Tools for controlling hardwares of Raspberry Pi

import (
	"fmt"
	"os/exec"
	"strings"
)

// Run given command with parameters and return combined output
func runCmd(cmdAndParams []string) (string, error) {
	if len(cmdAndParams) < 1 {
		return "", fmt.Errorf("No command provided")
	}

	output, err := exec.Command(cmdAndParams[0], cmdAndParams[1:]...).CombinedOutput()
	return strings.TrimRight(string(output), "\n"), err
}

// Reboot system
func RebootNow() (result string, err error) {
	return runCmd([]string{"sudo", "shutdown", "-r", "now"})
}

// Shutdown system
func ShutdownNow() (result string, err error) {
	return runCmd([]string{"sudo", "shutdown", "-h", "now"})
}

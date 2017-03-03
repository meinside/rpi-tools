package service

import (
	"fmt"
	"os/exec"
	"strings"
)

// Sudo run given command with parameters and return combined output
func sudoRunCmd(cmdAndParams []string) (string, error) {
	if len(cmdAndParams) < 1 {
		return "", fmt.Errorf("No command provided")
	}

	output, err := exec.Command("sudo", cmdAndParams...).CombinedOutput()
	return strings.TrimRight(string(output), "\n"), err
}

// Run `systemctl status is-active`
func SystemctlStatus(services []string) (statuses map[string]string, success bool) {
	statuses = make(map[string]string)

	args := []string{"systemctl", "is-active"}
	args = append(args, services...)

	output, _ := sudoRunCmd(args)
	for i, status := range strings.Split(output, "\n") {
		statuses[services[i]] = status
	}

	return statuses, true
}

// Run `systemctl start [service]`
func SystemctlStart(service string) (message string, success bool) {
	if output, err := sudoRunCmd([]string{"systemctl", "start", service}); err == nil {
		return output, true
	} else {
		return fmt.Sprintf("Failed to start service: %s", service), false
	}
}

// Run `systemctl stop [service]`
func SystemctlStop(service string) (message string, success bool) {
	if output, err := sudoRunCmd([]string{"systemctl", "stop", service}); err == nil {
		return output, true
	} else {
		return fmt.Sprintf("Failed to stop service: %s", service), false
	}
}

// Run `systemctl restart [service]`
func SystemctlRestart(service string) (message string, success bool) {
	if output, err := sudoRunCmd([]string{"systemctl", "restart", service}); err == nil {
		return output, true
	} else {
		return fmt.Sprintf("Failed to restart service: %s", service), false
	}
}

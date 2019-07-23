package service

import (
	"fmt"
	"strings"

	"github.com/meinside/rpi-tools/command"
)

// Run `systemctl status is-active`
func SystemctlStatus(services []string) (statuses map[string]string, success bool) {
	statuses = make(map[string]string)

	args := []string{"systemctl", "is-active"}
	args = append(args, services...)

	output, _ := command.SudoRun(args...)
	for i, status := range strings.Split(output, "\n") {
		statuses[services[i]] = status
	}

	return statuses, true
}

// Run `systemctl start [service]`
func SystemctlStart(service string) (message string, success bool) {
	if output, err := command.SudoRun("systemctl", "start", service); err == nil {
		return output, true
	} else {
		return fmt.Sprintf("Failed to start service: %s", service), false
	}
}

// Run `systemctl stop [service]`
func SystemctlStop(service string) (message string, success bool) {
	if output, err := command.SudoRun("systemctl", "stop", service); err == nil {
		return output, true
	} else {
		return fmt.Sprintf("Failed to stop service: %s", service), false
	}
}

// Run `systemctl restart [service]`
func SystemctlRestart(service string) (message string, success bool) {
	if output, err := command.SudoRun("systemctl", "restart", service); err == nil {
		return output, true
	} else {
		return fmt.Sprintf("Failed to restart service: %s", service), false
	}
}

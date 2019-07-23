package hardware

import (
	"github.com/meinside/rpi-tools/command"
)

// Tools for controlling hardwares of Raspberry Pi

// Reboot system
func RebootNow() (result string, err error) {
	return command.Run("sudo", "shutdown", "-r", "now")
}

// Shutdown system
func ShutdownNow() (result string, err error) {
	return command.Run("sudo", "shutdown", "-h", "now")
}

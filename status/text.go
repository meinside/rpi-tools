package status

// Tools for retrieving various statuses of Raspberry Pi
// with shell commands

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

// Get hostname
// (`hostname`)
func Hostname() (result string, err error) {
	return runCmd([]string{"hostname"})
}

// Get uname with '-a' parameter
// (`uname -a`)
func Uname() (result string, err error) {
	return runCmd([]string{"uname", "-a"})
}

// Get system uptime
// (`uptime`)
func Uptime() (result string, err error) {
	return runCmd([]string{"uptime"})
}

// Get disk usages
// (`df -h`)
func FreeSpaces() (result string, err error) {
	return runCmd([]string{"df", "-h"})
}

// Get memory split: arm and gpu
// (`vcgencmd get_mem arm; vcgencmd get_mem gpu`)
func MemorySplit() (result []string, err error) {
	var output string
	// arm memory
	output, err = runCmd([]string{"vcgencmd", "get_mem", "arm"})
	result = append(result, output)
	if err == nil {
		// gpu memory
		output, err = runCmd([]string{"vcgencmd", "get_mem", "gpu"})
		result = append(result, output)
	}
	return
}

// Get free memory
// (`free -o -h`)
func FreeMemory() (result string, err error) {
	return runCmd([]string{"free", "-o", "-h"})
}

// Get CPU temperature
// (`vcgencmd measure_temp`)
func CpuTemperature() (result string, err error) {
	return runCmd([]string{"vcgencmd", "measure_temp"})
}

// Get CPU information
// (`cat /proc/cpuinfo`)
func CpuInfo() (result string, err error) {
	return runCmd([]string{"cat", "/proc/cpuinfo"})
}

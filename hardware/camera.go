package hardware

// Tools for controlling various hardwares of Raspberry Pi

import (
	"fmt"
	"os/exec"
	"strconv"
)

const (
	// absolute path of raspistill
	RaspiStillBin = "/usr/bin/raspistill"
)

// capture an image with given width, height, and other parameters
// return the captured image's bytes
func CaptureRaspiStill(width, height int, cameraParams map[string]interface{}) (bytes []byte, err error) {
	// command line arguments
	args := []string{
		"-w", strconv.Itoa(width),
		"-h", strconv.Itoa(height),
		"-o", "-", // output to stdout
	}
	for k, v := range cameraParams {
		args = append(args, k)
		if v != nil {
			args = append(args, fmt.Sprintf("%v", v))
		}
	}

	// execute command
	if bytes, err := exec.Command(RaspiStillBin, args...).CombinedOutput(); err != nil {
		return []byte{}, err
	} else {
		return bytes, nil
	}
}

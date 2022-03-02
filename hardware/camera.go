package hardware

// Tools for controlling various hardwares of Raspberry Pi

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"time"
)

const (
	LibCameraStillBin               = "/usr/bin/libcamera-still"
	LibCameraStillRunTimeoutSeconds = 10
)

// CaptureStillImage captures an image with `libcamera-still`.
func CaptureStillImage(libcameraStillBinPath string, width, height int, cameraParams map[string]interface{}) (result []byte, err error) {
	// command line arguments
	args := []string{
		"--width", strconv.Itoa(width),
		"--height", strconv.Itoa(height),
		"--encoding", "jpg",
		"--output", "-", // output to stdout
	}
	for k, v := range cameraParams {
		args = append(args, k)
		if v != nil {
			args = append(args, fmt.Sprintf("%v", v))
		}
	}

	// execute command with timeout,
	cmd := exec.Command(libcameraStillBinPath, args...)
	var buffer bytes.Buffer
	cmd.Stdout = &buffer
	err = cmd.Start()
	if err == nil {
		done := make(chan error)
		go func() { done <- cmd.Wait() }()
		timeout := time.After(LibCameraStillRunTimeoutSeconds * time.Second)

		// and get its standard output
		select {
		case <-timeout:
			err = cmd.Process.Kill()
			if err == nil {
				err = fmt.Errorf("Command timed out: %s", libcameraStillBinPath)
			} else {
				err = fmt.Errorf("Command timed out, but failed to kill process: %s", libcameraStillBinPath)
			}
		case err = <-done:
			if err == nil {
				return buffer.Bytes(), nil
			} else {
				err = fmt.Errorf("Error running %s: %s", libcameraStillBinPath, err)
			}
		}
	}

	return nil, err
}

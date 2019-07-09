package status

// Tools for retrieving various statuses of Raspberry Pi

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Run given command with parameters and return combined output
func runCmd(cmdAndParams []string) (string, error) {
	if len(cmdAndParams) < 1 {
		return "", fmt.Errorf("No command provided")
	}

	output, err := exec.Command(cmdAndParams[0], cmdAndParams[1:]...).CombinedOutput()
	return strings.TrimRight(string(output), "\n"), err
}

// Hostname fetches hostname
// (`hostname`)
func Hostname() (result string, err error) {
	return runCmd([]string{"hostname"})
}

// Uname fetches uname with '-a' parameter
// (`uname -a`)
func Uname() (result string, err error) {
	return runCmd([]string{"uname", "-a"})
}

// Uptime fetches system uptime
// (`uptime`)
func Uptime() (result string, err error) {
	return runCmd([]string{"uptime"})
}

// FreeSpaces fetches disk usages
// (`df -h`)
func FreeSpaces() (result string, err error) {
	return runCmd([]string{"df", "-h"})
}

// MemorySplit fetches memory split: arm and gpu
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

// FreeMemory fetches free memory
// (`free -o -h`)
func FreeMemory() (result string, err error) {
	return runCmd([]string{"free", "-h"})
}

// MemoryUsage fetches system & heap allocated memory usage
func MemoryUsage() (sys, heap uint64) {
	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	return m.Sys, m.HeapAlloc
}

// CpuTemperature fetches CPU temperature
// (`vcgencmd measure_temp`)
func CpuTemperature() (result string, err error) {
	result, err = runCmd([]string{"vcgencmd", "measure_temp"})
	if err == nil {
		comps := strings.Split(result, "=") // eg: "temp=68.0'C"
		if len(comps) == 2 {
			return comps[1], nil
		}
	}
	return result, err
}

// CpuFrequency fetches frequency of arm clock
// (`vcgencmd measure_clock arm`)
func CpuFrequency() (result string, err error) {
	result, err = runCmd([]string{"vcgencmd", "measure_clock", "arm"})
	if err == nil {
		comps := strings.Split(result, "=") // eg: "frequency(48)=600169920"
		if len(comps) == 2 {
			num, _ := strconv.ParseFloat(strings.TrimSpace(comps[1]), 64)
			return fmt.Sprintf("%.1f MHz", num/1000.0/1000.0), nil
		}
	}
	return result, err
}

// CpuInfo fetches CPU information
// (`cat /proc/cpuinfo`)
func CpuInfo() (result string, err error) {
	return runCmd([]string{"cat", "/proc/cpuinfo"})
}

// IpAddresses fetches IP addresses
//
// http://play.golang.org/p/BDt3qEQ_2H
func IpAddresses() []string {
	ips := []string{}
	if ifaces, err := net.Interfaces(); err == nil {
		for _, iface := range ifaces {
			// skip
			if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
				continue
			}

			if addrs, err := iface.Addrs(); err == nil {
				for _, addr := range addrs {
					var ip net.IP
					switch v := addr.(type) {
					case *net.IPNet:
						ip = v.IP
					case *net.IPAddr:
						ip = v.IP
					}
					if ip == nil || ip.IsLoopback() {
						continue
					}
					ip = ip.To4()
					if ip == nil {
						continue
					}

					ips = append(ips, ip.String())
				}
			}
		}
	}

	return ips
}

// ExternalIpAddress fetches external IP address (https://gist.github.com/jniltinho/9788121)
func ExternalIpAddress() (ip string, err error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			IdleConnTimeout:       30 * time.Second,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	// http get request
	var req *http.Request
	if req, err = http.NewRequest("GET", "https://domains.google.com/checkip", nil); err == nil {
		// user-agent
		req.Header.Set("User-Agent", fmt.Sprintf("rpi-tools (golang; %s; %s)", runtime.GOOS, runtime.GOARCH))

		// http get
		var resp *http.Response
		resp, err = client.Do(req)

		if resp != nil {
			defer resp.Body.Close() // in case of http redirects
		}

		if err == nil && resp.StatusCode == 200 {
			var body []byte
			if body, err = ioutil.ReadAll(resp.Body); err == nil {
				ip := strings.TrimSpace(string(body))

				return ip, nil
			}

			err = fmt.Errorf("failed to read external ip: %s", err)
		} else {
			err = fmt.Errorf("failed to fetch external ip: %s (http %d)", err, resp.StatusCode)
		}
	}

	return "0.0.0.0", err
}

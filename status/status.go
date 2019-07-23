package status

// Tools for retrieving various statuses of Raspberry Pi

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/meinside/rpi-tools/command"
)

// Hostname fetches hostname
// (`hostname`)
func Hostname() (result string, err error) {
	return command.Run("hostname")
}

// Uname fetches uname with '-a' parameter
// (`uname -a`)
func Uname() (result string, err error) {
	return command.Run("uname", "-a")
}

// Uptime fetches system uptime
// (`uptime`)
func Uptime() (result string, err error) {
	return command.Run("uptime")
}

// FreeSpaces fetches disk usages
// (`df -h`)
func FreeSpaces() (result string, err error) {
	return command.Run("df", "-h")
}

// MemorySplit fetches memory split: arm and gpu
// (`vcgencmd get_mem arm; vcgencmd get_mem gpu`)
func MemorySplit() (result []string, err error) {
	var output string
	// arm memory
	output, err = command.Run("vcgencmd", "get_mem", "arm")
	result = append(result, output)
	if err == nil {
		// gpu memory
		output, err = command.Run("vcgencmd", "get_mem", "gpu")
		result = append(result, output)
	}
	return
}

// FreeMemory fetches free memory
// (`free -o -h`)
func FreeMemory() (result string, err error) {
	return command.Run("free", "-h")
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
	result, err = command.Run("vcgencmd", "measure_temp")
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
	result, err = command.Run("vcgencmd", "measure_clock", "arm")
	if err == nil {
		comps := strings.Split(result, "=") // eg: "frequency(48)=600169920"
		if len(comps) == 2 {
			num, _ := strconv.ParseFloat(strings.TrimSpace(comps[1]), 64)
			return fmt.Sprintf("%.1f MHz", num/1000.0/1000.0), nil
		}
	}
	return result, err
}

// CpuThrottled returns whether the system is throttled or not
func CpuThrottled() (result string, err error) {
	result, err = command.Run("vcgencmd", "get_throttled")
	if err == nil {
		comps := strings.Split(result, "=") // eg: throttled=0x50000
		if len(comps) == 2 {
			num, _ := strconv.ParseInt(strings.Replace(strings.TrimSpace(comps[1]), "0x", "", -1), 16, 64)

			results := []string{}

			// https://www.raspberrypi.org/forums/viewtopic.php?f=63&t=147781&start=50#p972790
			if num&1 > 0 {
				// under-voltage
				results = append(results, "under-voltage now")
			}
			if num&(1<<1) > 0 {
				// arm frequency capped
				results = append(results, "arm freq capped now")
			}
			if num&(1<<2) > 0 {
				// currently throttled
				results = append(results, "throttled now")
			}
			if num&(1<<16) > 0 && num&1 <= 0 {
				// under-voltage has occurred
				results = append(results, "under-voltage before")
			}
			if num&(1<<17) > 0 && num&(1<<1) <= 0 {
				// arm frequency capped has occurred
				results = append(results, "arm freq capped before")
			}
			if num&(1<<18) > 0 && num&(1<<2) <= 0 {
				// throttling has occurred
				results = append(results, "throttled before")
			}

			if len(results) <= 0 {
				result = "ok"
			} else {
				result = strings.Join(results, ", ")
			}

			return result, nil
		}
	}
	return result, err
}

// CpuInfo fetches CPU information
// (`cat /proc/cpuinfo`)
func CpuInfo() (result string, err error) {
	return command.Run("cat", "/proc/cpuinfo")
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

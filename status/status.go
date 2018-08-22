package status

// Tools for retrieving various statuses of Raspberry Pi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
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
	return runCmd([]string{"vcgencmd", "measure_temp"})
}

// CpuInfo fetches CPU information
// (`cat /proc/cpuinfo`)
func CpuInfo() (result string, err error) {
	return runCmd([]string{"cat", "/proc/cpuinfo"})
}

// Free Geo IP information provided by http://geoip.nekudo.com/

// CityValue is a struct for city value
type CityValue interface{} // XXX - can be a string or a bool value

// GeoInfo struct
type GeoInfo struct {
	City     CityValue     `json:"city"`
	Country  GeoIpCountry  `json:"country"`
	Location GeoIpLocation `json:"location"`
	Ip       string        `json:"ip"`
}

// GeoIpCountry struct
type GeoIpCountry struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// GeoIpLocation struct
type GeoIpLocation struct {
	AccuracyRadius int     `json:"accuracy_radius"`
	Latitude       float32 `json:"latitude"`
	Longitude      float32 `json:"longitude"`
	Timezone       string  `json:"time_zone"`
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
func ExternalIpAddress() (string, error) {
	var resp *http.Response
	var err error

	resp, err = http.Get("http://icanhazip.com")

	if resp != nil {
		defer resp.Body.Close() // in case of http redirects
	}

	if err == nil && resp.StatusCode == 200 {
		var body []byte
		if body, err = ioutil.ReadAll(resp.Body); err == nil {
			return strings.TrimSpace(string(body)), nil
		}

		log.Printf("Failed to read external ip: %s\n", err)
	} else {
		log.Printf("Failed to fetch external ip: %s (http %d)\n", err, resp.StatusCode)
	}

	return "0.0.0.0", err
}

// GeoLocation fetches GeoInfo result with given IP address
func GeoLocation(ip string) (GeoInfo, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 300 * time.Second,
			}).Dial,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	var req *http.Request
	var resp *http.Response
	var err error
	if req, err = http.NewRequest("GET", "https://geoip.nekudo.com/api/"+ip, nil); err == nil {
		resp, err = client.Do(req)

		if resp != nil {
			defer resp.Body.Close() // in case of http redirects
		}

		if err == nil {
			var body []byte
			if body, err = ioutil.ReadAll(resp.Body); err == nil {
				if resp.StatusCode == 200 {
					var jsonResp GeoInfo
					if err = json.Unmarshal(body, &jsonResp); err == nil {
						return jsonResp, nil
					}

					log.Printf("Failed to parse geo info json: %s\n", err)
				} else {
					log.Printf("Geo info HTTP error %d\n", resp.StatusCode)
				}
			} else {
				log.Printf("Failed to read geo info response: %s\n", err)
			}
		} else {
			log.Printf("Failed to request geo info: %s\n", err)
		}
	} else {
		log.Printf("Failed to create a geo info request: %s\n", err)
	}

	return GeoInfo{}, err
}

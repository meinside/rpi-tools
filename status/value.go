package status

// Tools for retrieving various statuses of Raspberry Pi

import (
	"net"
	"runtime"
)

// Get system & heap allocated memory usage
func MemoryUsage() (sys, heap uint64) {
	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	return m.Sys, m.HeapAlloc
}

// Get IP addresses
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

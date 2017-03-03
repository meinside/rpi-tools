# rpi-tools

Various tools for Raspberry Pi, written in Golang.

## How to get

```bash
$ go get -u github.com/meinside/rpi-tools/...
```

## Usage (example)

### Status

```go
package main

import (
	"fmt"
	"strings"

	"github.com/meinside/rpi-tools/status"
)

func main() {
	printStatuses1()
	printStatuses2()
}

func printStatuses1() {
	if result, err := status.Hostname(); err == nil {
		fmt.Printf("> hostname\n%s\n\n", result)
	} else {
		fmt.Printf("* error: %s\n", err)
	}

	if result, err := status.Uname(); err == nil {
		fmt.Printf("> uname\n%s\n\n", result)
	} else {
		fmt.Printf("* error: %s\n", err)
	}

	if result, err := status.Uptime(); err == nil {
		fmt.Printf("> uptime\n%s\n\n", result)
	} else {
		fmt.Printf("* error: %s\n", err)
	}

	if result, err := status.FreeSpaces(); err == nil {
		fmt.Printf("> free spaces\n%s\n\n", result)
	} else {
		fmt.Printf("* error: %s\n", err)
	}

	if result, err := status.MemorySplit(); err == nil {
		fmt.Printf("> memory split\n%s\n\n", strings.Join(result, ", "))
	} else {
		fmt.Printf("* error: %s\n", err)
	}

	if result, err := status.FreeMemory(); err == nil {
		fmt.Printf("> free memory\n%s\n\n", result)
	} else {
		fmt.Printf("* error: %s\n", err)
	}

	if result, err := status.CpuTemperature(); err == nil {
		fmt.Printf("> cpu temp\n%s\n\n", result)
	} else {
		fmt.Printf("* error: %s\n", err)
	}

	if result, err := status.CpuInfo(); err == nil {
		fmt.Printf("> cpu info\n%s\n\n", result)
	} else {
		fmt.Printf("* error: %s\n", err)
	}
}

func printStatuses2() {
	sys, heap := status.MemoryUsage()
	fmt.Printf("> memory usage\nsystem: %.1fMB, heap: %.1fMB\n\n",
		float32(sys)/1024/1024,
		float32(heap)/1024/1024,
	)

	ips := status.IpAddresses()
	fmt.Printf("> assigned ips\n%s\n\n", strings.Join(ips, ", "))
}
```

It will print out:

```
> hostname
raspberry

> uname
Linux raspberry 4.4.48-v7+ #964 SMP Mon Feb 13 16:57:51 GMT 2017 armv7l GNU/Linux

> uptime
 15:06:05 up  5:03,  1 user,  load average: 0.30, 0.17, 0.11

> free spaces
Filesystem      Size  Used Avail Use% Mounted on
/dev/root        20G  4.4G   15G  24% /
devtmpfs        483M     0  483M   0% /dev
tmpfs           487M  4.0K  487M   1% /dev/shm
tmpfs           487M   19M  468M   4% /run
tmpfs           5.0M  4.0K  5.0M   1% /run/lock
tmpfs           487M     0  487M   0% /sys/fs/cgroup
/dev/mmcblk0p1   60M   21M   39M  36% /boot
/dev/sda3       438G  342G   74G  83% /home

> memory split
arm=992M, gpu=16M

> free memory
             total       used       free     shared    buffers     cached
Mem:          973M       816M       156M        29M        80M       463M
Swap:         1.1G         0B       1.1G

> cpu temp
temp=52.6'C

> cpu info
processor       : 0
model name      : ARMv7 Processor rev 4 (v7l)
BogoMIPS        : 76.80
Features        : half thumb fastmult vfp edsp neon vfpv3 tls vfpv4 idiva idivt vfpd32 lpae evtstrm crc32
CPU implementer : 0x41
CPU architecture: 7
CPU variant     : 0x0
CPU part        : 0xd03
CPU revision    : 4

processor       : 1
model name      : ARMv7 Processor rev 4 (v7l)
BogoMIPS        : 76.80
Features        : half thumb fastmult vfp edsp neon vfpv3 tls vfpv4 idiva idivt vfpd32 lpae evtstrm crc32
CPU implementer : 0x41
CPU architecture: 7
CPU variant     : 0x0
CPU part        : 0xd03
CPU revision    : 4

processor       : 2
model name      : ARMv7 Processor rev 4 (v7l)
BogoMIPS        : 76.80
Features        : half thumb fastmult vfp edsp neon vfpv3 tls vfpv4 idiva idivt vfpd32 lpae evtstrm crc32
CPU implementer : 0x41
CPU architecture: 7
CPU variant     : 0x0
CPU part        : 0xd03
CPU revision    : 4

processor       : 3
model name      : ARMv7 Processor rev 4 (v7l)
BogoMIPS        : 76.80
Features        : half thumb fastmult vfp edsp neon vfpv3 tls vfpv4 idiva idivt vfpd32 lpae evtstrm crc32
CPU implementer : 0x41
CPU architecture: 7
CPU variant     : 0x0
CPU part        : 0xd03
CPU revision    : 4

Hardware        : BCM2709
Revision        : a02082
Serial          : 0000000000000000

> memory usage
system: 18.3MB, heap: 0.1MB

> assigned ips
192.168.0.7

```

### Service (systemd)

```go
package main

import (
	"fmt"

	"github.com/meinside/rpi-tools/service"
)

func main() {
	testSystemd()
}

func testSystemd() {
	var msg string
	var success bool

	var statuses map[string]string
	statuses, success = service.SystemctlStatus([]string{"some-service1", "some-service2"})
	fmt.Printf("> service statuses:\n")
	for service, status := range statuses {
		fmt.Printf("  %s: %s\n", service, status)
	}

	msg, success = service.SystemctlStart("some-service3")
	if success {
		fmt.Printf("> %s\n", msg)
	} else {
		fmt.Printf("* %s\n", msg)
	}

	msg, success = service.SystemctlStop("some-service3")
	if success {
		fmt.Printf("> %s\n", msg)
	} else {
		fmt.Printf("* %s\n", msg)
	}
}
```

## Todos

- [ ] Add some more useful functions

## License

MIT

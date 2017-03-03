package status

// Tools for retrieving various statuses of Raspberry Pi

import (
	"runtime"
)

// Get system & heap allocated memory usage
func MemoryUsage() (sys, heap uint64) {
	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	return m.Sys, m.HeapAlloc
}

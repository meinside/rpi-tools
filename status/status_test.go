package status

import (
	"testing"
)

func TestIpAddresses(t *testing.T) {
	ips := IpAddresses()
	if len(ips) <= 0 {
		t.Errorf("failed to get local ip addresses")
	}
}

func TestExternalIpAddress(t *testing.T) {
	_, err := ExternalIpAddress()
	if err != nil {
		t.Errorf("failed to get external ip addresses: %s", err)
	}
}

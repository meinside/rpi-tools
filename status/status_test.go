package status

import (
	"testing"
)

func TestIPAddresses(t *testing.T) {
	ips := IPAddresses()
	if len(ips) <= 0 {
		t.Errorf("failed to get local ip addresses")
	}
}

func TestExternalIpAddress(t *testing.T) {
	_, err := ExternalIPAddress()
	if err != nil {
		t.Errorf("failed to get external ip addresses: %s", err)
	}
}

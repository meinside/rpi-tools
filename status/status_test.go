package status

import (
	"testing"
)

// If this test fails, visit http://geoip.nekudo.com/ and check for any change.
func TestGeoLocation(t *testing.T) {
	ip := "8.8.8.8"
	var lat float32 = 37.751
	var lon float32 = -97.822

	if geoInfo, err := GeoLocation(ip); err == nil {
		if geoInfo.Ip != ip {
			t.Errorf("returned ip differs from the requested one: %s", geoInfo.Ip)
		}
		if geoInfo.Location.Latitude != lat || geoInfo.Location.Longitude != lon {
			t.Errorf("returned location seems to be different from the expected one: (%.3f, %.3f)", geoInfo.Location.Latitude, geoInfo.Location.Longitude)
		}
	} else {
		t.Errorf("failed to get geo location: %s\n", err)
	}
}

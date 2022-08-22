package utils

import "testing"

func TestGeoHashEncode(t *testing.T) {
	encode, err := GeoHashEncode(104.05503, 30.562251, 6)
	if err != nil {
		t.Errorf("Geo hash encode err: %v.", err)
		return
	}
	t.Log(encode)
}

func TestGeoHashDecode(t *testing.T) {
	longitude, latitude, err := GeoHashDecode("wm3vzg")
	if err != nil {
		t.Errorf("Geo hash decode err: %v.", err)
		return
	}
	t.Log(longitude, latitude)
}

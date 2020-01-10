package xmp

import (
	"log"
	"testing"
	"time"
)

func TestWithoutAlpha(t *testing.T) {
	profileInput := Profile{
		Name:      "Toby Christopher",
		Timestamp: time.Now().UnixNano(),
		Lat:       53.5395,
		Long:      10.0051,
		Location:  "HafenCity, Hamburg, Germany",
	}

	err := WriteXMP("samples/webp_noalpha.webp", profileInput.Name, profileInput.Timestamp, profileInput.Lat, profileInput.Long, profileInput.Location)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	profile, err := ReadXMP("samples/webp_noalpha.webp")
	if err != nil {
		t.Errorf("error: %v", err)
	}
	if *profile != profileInput {
		t.Errorf("error:\nexpect %v\nbut received %v", profileInput, *profile)
	}
	// Test passed
	log.Printf("profile: %v", profile)
}

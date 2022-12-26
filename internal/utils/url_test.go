package utils

import "testing"

func TestGetBaseUrl(t *testing.T) {
	url := "http://localhost:8080/stream/frag0.ts"

	baseUrl, err := GetBaseUrl(url)
	if err != nil {
		t.Error(err)
	}

	if baseUrl != "http://localhost:8080/stream/" {
		t.Errorf("Incorrect base url, %s", baseUrl)
	}
}

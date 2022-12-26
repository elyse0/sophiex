package utils

import (
	"regexp"
)

func GetBaseUrl(url string) (string, error) {
	re := regexp.MustCompile("https?://[^?#]+/")
	match := re.FindStringSubmatch(url)

	return match[0], nil
}

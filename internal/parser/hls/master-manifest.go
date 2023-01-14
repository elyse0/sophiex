package hls

import (
	"fmt"
	"strings"
)

type MasterManifest struct {
	Manifest string
}

type Stream struct {
	Bandwidth  int64
	Resolution struct {
		Height int64
		Width  int64
	}
	Codecs string
}

func (masterManifest MasterManifest) Parse() (string, error) {
	for _, line := range strings.Split(strings.TrimSuffix(masterManifest.Manifest, "\n"), "\n") {
		fmt.Println(line)
	}

	return "", nil
}

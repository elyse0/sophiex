package parser

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type HlsFragment struct {
	MediaSequence int     `json:"mediaSequence"`
	Discontinuity int     `json:"discontinuity"`
	Url           string  `json:"url"`
	Duration      float64 `json:"duration"`
	Start         float32 `json:"start"`
	End           float32 `json:"end"`
}

type HlsMediaManifest struct {
	manifest    string
	manifestUrl string
}

func (mediaManifest HlsMediaManifest) GetFragments() ([]HlsFragment, error) {
	var fragments []HlsFragment

	mediaSequence := 0
	discontinuity := 0
	var duration float64 = 0

	for _, line := range strings.Split(strings.TrimSuffix(mediaManifest.manifest, "\n"), "\n") {
		line = strings.TrimSpace(line)

		if !strings.HasPrefix(line, "#") {
			var fragmentUrl string

			urlMatch, _ := regexp.MatchString("^https?://", line)
			if urlMatch {
				fragmentUrl = line
			} else {
				urlJoin, _ := url.JoinPath(mediaManifest.manifestUrl, line)
				fragmentUrl = urlJoin
			}

			fragments = append(fragments, HlsFragment{
				MediaSequence: mediaSequence,
				Discontinuity: discontinuity,
				Url:           fragmentUrl,
				Duration:      duration,
				Start:         0,
				End:           0,
			})

			mediaSequence += 1
		} else if strings.HasPrefix(line, "#EXT-X-MAP") {

		} else if strings.HasPrefix(line, "#EXTINF") {
			re := regexp.MustCompile(`#EXTINF\s*:\s*([\d+.]+)`)
			match := re.FindStringSubmatch(line)
			duration, _ = strconv.ParseFloat(match[1], 64)
		} else if strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE") {
			re := regexp.MustCompile(`#EXT-X-MEDIA-SEQUENCE\s*:\s*(\d+)`)
			match := re.FindStringSubmatch(line)
			mediaSequence, _ = strconv.Atoi(match[1])
		}
	}

	return fragments, nil
}

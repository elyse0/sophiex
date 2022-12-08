package parser

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type HlsFragment struct {
	mediaSequence int
	discontinuity int
	url           string
	duration      float32
	start         float32
	end           float32
}

type HlsMediaManifest struct {
	manifest    string
	manifestUrl string
}

func (mediaManifest HlsMediaManifest) GetFragments() {
	var fragments []HlsFragment

	mediaSequence := 0
	discontinuity := 0
	var duration float32 = 0

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
				mediaSequence: mediaSequence,
				discontinuity: discontinuity,
				url:           fragmentUrl,
				duration:      duration,
				start:         0,
				end:           0,
			})

			mediaSequence += 1
		} else if strings.HasPrefix("#EXT-X-MAP", line) {

		} else if strings.HasPrefix("#EXT-X-MEDIA-SEQUENCE", line) {
			re := regexp.MustCompile(`#EXT-X-MEDIA-SEQUENCE\s*:\s*(\d+)`)
			match := re.FindStringSubmatch(line)
			mediaSequence, _ = strconv.Atoi(match[1])
		}
	}
}

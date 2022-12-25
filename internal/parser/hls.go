package parser

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type HlsFragment struct {
	MediaSequence int    `json:"mediaSequence"`
	Discontinuity int    `json:"discontinuity"`
	Url           string `json:"url"`
	Duration      int64  `json:"duration"` // Milliseconds
	Start         int64  `json:"start"`    // Unix milliseconds
	End           int64  `json:"end"`      // Unix milliseconds
}

type HlsMediaManifest struct {
	manifest     string
	manifestUrl  string
	isLivestream bool
}

func (mediaManifest HlsMediaManifest) GetFragments() ([]HlsFragment, error) {
	var fragments []HlsFragment

	var programDateTime int64

	// Some livestreams don't provide program date times, this way at least we can approximate it.
	if mediaManifest.isLivestream {
		programDateTime = time.Now().UnixMilli()
	} else {
		programDateTime = 0
	}

	mediaSequence := 0
	discontinuity := 0
	var duration int64 = 0

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

			fragmentStart := programDateTime
			fragmentEnd := programDateTime + duration
			programDateTime += duration

			fragments = append(fragments, HlsFragment{
				MediaSequence: mediaSequence,
				Discontinuity: discontinuity,
				Url:           fragmentUrl,
				Duration:      duration,
				Start:         fragmentStart,
				End:           fragmentEnd,
			})

			mediaSequence += 1
		} else if strings.HasPrefix(line, "#EXT-X-MAP") {

		} else if strings.HasPrefix(line, "#EXTINF") {
			re := regexp.MustCompile(`#EXTINF\s*:\s*([\d+.]+)`)
			match := re.FindStringSubmatch(line)
			durationFloat, _ := strconv.ParseFloat(match[1], 64)
			duration = int64(durationFloat * 1000)
		} else if strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE") {
			re := regexp.MustCompile(`#EXT-X-MEDIA-SEQUENCE\s*:\s*(\d+)`)
			match := re.FindStringSubmatch(line)
			mediaSequence, _ = strconv.Atoi(match[1])
		} else if strings.HasPrefix(line, "#EXT-X-PROGRAM-DATE-TIME") {
			re := regexp.MustCompile(`#EXT-X-PROGRAM-DATE-TIME\s*:\s*([\w-:]+)`)
			match := re.FindStringSubmatch(line)
			tt, _ := time.Parse(time.RFC3339, match[1])
			programDateTime = tt.UnixMilli()
		}
	}

	return fragments, nil
}

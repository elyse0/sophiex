package parser

import (
	"fmt"
	"regexp"
	"sophiex/internal/fragment"
	"sophiex/internal/utils"
	"strconv"
	"strings"
	"time"
)

type HlsInitialization struct {
	Url       string             `json:"url"`
	ByteRange fragment.ByteRange `json:"byteRange"`
}

func (initialization *HlsInitialization) IsEmpty() bool {
	if initialization.Url == "" {
		return true
	}

	return false
}

type HlsFragment struct {
	MediaSequence int                 `json:"mediaSequence"`
	Discontinuity int                 `json:"discontinuity"`
	Url           string              `json:"url"`
	Duration      int64               `json:"duration"`   // Milliseconds
	Decryption    fragment.Decryption `json:"decryption"` // Milliseconds
	ByteRange     fragment.ByteRange  `json:"byteRange"`
	Start         int64               `json:"start"` // Unix milliseconds
	End           int64               `json:"end"`   // Unix milliseconds
}

func (fragment *HlsFragment) IsEmpty() bool {
	if fragment.Url == "" {
		return true
	}

	return false
}

type HlsMediaManifest struct {
	Manifest     string
	ManifestUrl  string
	IsLivestream bool
}

type HlsMediaManifestParseResult struct {
	Initialization HlsInitialization
	Fragments      []HlsFragment
}

func (mediaManifest HlsMediaManifest) Parse() (HlsMediaManifestParseResult, error) {
	var initialization HlsInitialization
	var fragments []HlsFragment

	var programDateTime int64

	// Some livestreams don't provide program date times, this way at least we can approximate it.
	if mediaManifest.IsLivestream {
		programDateTime = time.Now().UnixMilli()
	} else {
		programDateTime = 0
	}

	mediaSequence := 0
	discontinuity := 0
	var duration int64 = 0
	decryption := fragment.Decryption{}
	byteRange := fragment.ByteRange{}

	for _, line := range strings.Split(strings.TrimSuffix(mediaManifest.Manifest, "\n"), "\n") {
		line = strings.TrimSpace(line)

		if !strings.HasPrefix(line, "#") {
			var fragmentUrl string

			urlMatch, _ := regexp.MatchString("^https?://", line)
			if urlMatch {
				fragmentUrl = line
			} else {
				if strings.HasPrefix(line, "/") {
					re := regexp.MustCompile("https?://[^/]+")
					match := re.FindStringSubmatch(mediaManifest.ManifestUrl)
					fragmentUrl = match[0] + line
				} else {
					baseUrl, _ := utils.GetBaseUrl(mediaManifest.ManifestUrl)
					fragmentUrl = baseUrl + line
				}
			}

			fragmentStart := programDateTime
			fragmentEnd := programDateTime + duration
			programDateTime += duration

			fragments = append(fragments, HlsFragment{
				MediaSequence: mediaSequence,
				Discontinuity: discontinuity,
				Url:           fragmentUrl,
				Duration:      duration,
				Decryption:    decryption,
				ByteRange:     byteRange,
				Start:         fragmentStart,
				End:           fragmentEnd,
			})

			mediaSequence += 1
		} else if strings.HasPrefix(line, "#EXT-X-MAP") {
			re := regexp.MustCompile(`#EXT-X-MAP\s*:\s*URI\s*=\s*"([^"]+)"(?:,BYTERANGE="(\d+)@(\d+))?`)
			match := re.FindStringSubmatch(line)

			var initializationUrl string
			urlMatch, _ := regexp.MatchString("^https?://", match[1])
			if urlMatch {
				initializationUrl = match[1]
			} else {
				baseUrl, _ := utils.GetBaseUrl(mediaManifest.ManifestUrl)
				initializationUrl = baseUrl + match[1]
			}

			byteRange := fragment.ByteRange{}
			if len(match) > 2 {
				length, _ := strconv.Atoi(match[2])
				start, _ := strconv.Atoi(match[3])

				byteRange = fragment.ByteRange{
					Start: start,
					End:   start + length,
				}
			}

			initialization = HlsInitialization{
				Url:       initializationUrl,
				ByteRange: byteRange,
			}

		} else if strings.HasPrefix(line, "#EXTINF") {
			re := regexp.MustCompile(`#EXTINF\s*:\s*([\d+.]+)`)
			match := re.FindStringSubmatch(line)
			durationFloat, _ := strconv.ParseFloat(match[1], 64)
			duration = int64(durationFloat * 1000)
			byteRange = fragment.ByteRange{}
		} else if strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE") {
			re := regexp.MustCompile(`#EXT-X-MEDIA-SEQUENCE\s*:\s*(\d+)`)
			match := re.FindStringSubmatch(line)
			mediaSequence, _ = strconv.Atoi(match[1])
		} else if strings.HasPrefix(line, "#EXT-X-PROGRAM-DATE-TIME") {
			re := regexp.MustCompile(`#EXT-X-PROGRAM-DATE-TIME\s*:\s*([\w-:]+)`)
			match := re.FindStringSubmatch(line)
			tt, _ := time.Parse(time.RFC3339, match[1])
			programDateTime = tt.UnixMilli()
		} else if strings.HasPrefix(line, "#EXT-X-KEY") {
			re := regexp.MustCompile(`#EXT-X-KEY\s*:\s*METHOD\s*=\s*NONE`)
			match := re.FindStringSubmatch(line)
			if len(match) != 0 {
				decryption = fragment.Decryption{}
				continue
			}

			re = regexp.MustCompile(`#EXT-X-KEY\s*:\s*METHOD\s*=\s*AES-128\s*,URI\s*=\s*"([^"]+)`)
			match = re.FindStringSubmatch(line)
			if len(match) != 0 {
				decryption = fragment.Decryption{
					Method: "AES-128",
					Uri:    match[1],
				}
			}
		} else if strings.HasPrefix(line, "#EXT-X-BYTERANGE") {
			re := regexp.MustCompile(`#EXT-X-BYTERANGE\s*:\s*(\d+)(?:@(\d+))?`)
			match := re.FindStringSubmatch(line)
			length, _ := strconv.Atoi(match[1])
			start, _ := strconv.Atoi(match[2])

			byteRange = fragment.ByteRange{
				Start: start,
				End:   start + length,
			}
		}
	}

	return HlsMediaManifestParseResult{
		Initialization: initialization,
		Fragments:      fragments,
	}, nil
}

type HlsMasterManifest struct {
	Manifest string
}

type HlsStream struct {
	Bandwidth  int64
	Resolution struct {
		Height int64
		Width  int64
	}
	Codecs string
}

func (masterManifest HlsMasterManifest) Parse() (string, error) {
	for _, line := range strings.Split(strings.TrimSuffix(masterManifest.Manifest, "\n"), "\n") {
		fmt.Println(line)
	}

	return "", nil
}

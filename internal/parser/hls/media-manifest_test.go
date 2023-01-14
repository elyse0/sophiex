package hls

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sophiex/internal/fragment"
	"testing"
)

func TestHlsMediaManifest_GetFragments(t *testing.T) {
	fileContent, err := os.ReadFile("samples/simple-hls-media-manifest.m3u8")
	if err != nil {
		t.Error(err)
	}

	mediaManifest := MediaManifest{
		Manifest:     string(fileContent),
		ManifestUrl:  "http://localhost:8080/stream_0/stream_0.m3u8",
		IsLivestream: false,
	}

	parseResult, err := mediaManifest.Parse()
	if err != nil {
		t.Error(err)
	}

	initializationJson, err := json.MarshalIndent(parseResult.Initialization, "", "\t")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(initializationJson))

	fragmentsJson, err := json.MarshalIndent(parseResult.Fragments, "", "\t")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(fragmentsJson))

	expectedFragments := []Fragment{
		{
			Generic: fragment.Generic{
				Url: "http://localhost:8080/stream_0/data00.ts?test=test",
			},
			MediaSequence: 0,
			Discontinuity: 0,
			Duration:      3200,
			Start:         0,
			End:           3200,
		},
		{
			Generic: fragment.Generic{
				Url: "http://localhost:8080/stream_0/data01.ts",
			},
			MediaSequence: 1,
			Discontinuity: 0,
			Duration:      1600,
			Start:         3200,
			End:           4800,
		},
		{
			Generic: fragment.Generic{
				Url: "http://localhost:8080/stream_0/data02.ts",
			},
			MediaSequence: 2,
			Discontinuity: 0,
			Duration:      1600,
			Start:         4800,
			End:           6400,
		},
	}

	if !reflect.DeepEqual(parseResult.Fragments, expectedFragments) {
		t.Errorf("Unexpected fragments")
	}
}

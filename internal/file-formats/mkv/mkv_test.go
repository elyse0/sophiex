package mkv

import (
	"encoding/json"
	"fmt"
	"os"
	"sophiex/internal/file-formats/ebml"
	"testing"
)

func TestHeader(t *testing.T) {
	file, err := os.Open("assets/big-buck-bunny-720p-1mb.mkv")
	if err != nil {
		t.Error(err)
	}

	document, err := ebml.Parse(file)
	if err != nil {
		t.Error(err)
	}

	documentJson, err := json.MarshalIndent(document, "", " ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(documentJson))
}

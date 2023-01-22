package variable_size_integer

import (
	"bytes"
	"encoding/hex"
	"testing"
)

// 4.4 VINT Examples, https://datatracker.ietf.org/doc/html/rfc8794#section-4.4
func TestMultiByteVariableSizeInteger(t *testing.T) {
	hexValues := []string{"82", "4002", "200002", "10000002"}
	for _, hexValue := range hexValues {
		variableSizeInteger, err := hex.DecodeString(hexValue)
		if err != nil {
			t.Error(err)
		}

		reader := bytes.NewReader(variableSizeInteger)

		integer, err := Read(reader)
		if err != nil {
			t.Error(err)
		}

		if integer != 0x02 {
			t.Errorf("Integer %d is not equal to 0x02", integer)
		}
	}
}

package ebml

import (
	"encoding/binary"
	"fmt"
	"io"
	"sophiex/internal/file-formats/ebml/components"
	variableSizeInteger "sophiex/internal/file-formats/ebml/variable-size-integer"
)

// Document https://datatracker.ietf.org/doc/html/rfc8794#section-8
type Document struct {
	Header *components.Header `json:"header"`
	Body   *components.Body   `json:"body"`
}

func parseString(reader io.Reader) (*string, error) {
	return nil, nil
}

// https://datatracker.ietf.org/doc/html/rfc8794#name-unsigned-integer-element
// An Unsigned Integer Element stores an integer (meaning that it can be written without a
// fractional component) that could be positive or zero.
func readUnsignedInteger(reader io.Reader) (uint64, error) {
	// An Unsigned Integer Element MUST declare a length from zero to eight octets.
	integerLength, err := variableSizeInteger.Read(reader)
	if err != nil {
		return 0, err
	}

	// If the EBML Element is not defined to have a default value, then an Unsigned Integer Element
	// with a zero-octet length represents an integer value of zero.

	integerBytes := make([]byte, integerLength)
	_, err = reader.Read(integerBytes)
	if err != nil {
		return 0, err
	}

	// Because EBML limits Unsigned Integers to 8 octets in length, an Unsigned Integer Element
	// stores a number from 0 to 18,446,744,073,709,551,615.
	if integerLength < 8 {
		integerBytes = append(make([]byte, 8-integerLength), integerBytes...)
	}

	return binary.BigEndian.Uint64(integerBytes), nil
}

func parseHeader(reader io.Reader) (*components.Header, error) {
	headerSize, err := variableSizeInteger.Read(reader)
	if err != nil {
		return nil, err
	}

	fmt.Println("HeaderSize: ", headerSize)

	versionElement, err := variableSizeInteger.ReadWithMarker(reader)
	if err != nil {
		return nil, err
	}

	if versionElement != EBMLVersion {
		return nil, fmt.Errorf("expected EBML Version %x, got %x", EBMLElement, versionElement)
	}

	version, err := readUnsignedInteger(reader)
	if err != nil {
		return nil, err
	}

	return &components.Header{
		DocType:        "Header",
		DocTypeVersion: version,
	}, nil
}

func Parse(reader io.Reader) (*Document, error) {
	elementId, err := variableSizeInteger.ReadWithMarker(reader)
	if err != nil {
		return nil, err
	}

	if elementId != EBMLElement {
		return nil, fmt.Errorf("expected EBML Element %x, got %x", EBMLElement, elementId)
	}

	header, err := parseHeader(reader)
	if err != nil {
		return nil, err
	}

	return &Document{
		Header: header,
	}, err
}

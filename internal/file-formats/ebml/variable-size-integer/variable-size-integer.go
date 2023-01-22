package variable_size_integer

import (
	"encoding/binary"
	"io"
	"math/bits"
)

// 4.1. VINT_WIDTH, https://datatracker.ietf.org/doc/html/rfc8794#section-4.1
func getWidth(firstByte byte) int {
	// Each Variable-Size Integer starts with a VINT_WIDTH followed by a VINT_MARKER.

	// VINT_WIDTH is a sequence of zero or more bits of value 0 and is terminated by the VINT_MARKER,
	// which is a single bit of value 1. The total length in bits of both VINT_WIDTH and VINT_MARKER
	// is the total length in octets in of the Variable-Size Integer.
	return bits.LeadingZeros8(firstByte) + 1
}

// Read
// 4. Variable-Size Integer, https://datatracker.ietf.org/doc/html/rfc8794#section-4.
func read(reader io.Reader, removeMarker bool) (uint64, error) {
	// The Variable-Size Integer is composed of a VINT_WIDTH, VINT_MARKER, and VINT_DATA, in that order.
	firstOctet := make([]byte, 1)
	_, err := reader.Read(firstOctet)
	if err != nil {
		return 0, err
	}

	width := getWidth(firstOctet[0])

	// The bits required for the VINT_WIDTH and the VINT_MARKER use one out of every eight bits
	// of the total length of the Variable-Size Integer. Thus, a Variable-Size Integer of 1-octet length
	// supplies 7 bits for VINT_DATA, a 2-octet length supplies 14 bits for VINT_DATA,
	// and a 3-octet length supplies 21 bits for VINT_DATA.
	if removeMarker {
		firstOctet[0] = firstOctet[0] & (0b11111111 >> width)
	}

	// The VINT_DATA portion of the Variable-Size Integer includes all data following (but not including)
	// the VINT_MARKER until end of the Variable-Size Integer whose length is derived from the VINT_WIDTH.

	// If the number of bits required for VINT_DATA is less than the bit size of VINT_DATA,
	// then VINT_DATA MUST be zero-padded to the left to a size that fits.
	var data []byte
	if width == 1 {
		data = make([]byte, 8)
		data[len(data)-1] = firstOctet[0]
	} else {
		data = make([]byte, width-1)
		_, err = reader.Read(data)
		if err != nil {
			return 0, err
		}

		data = append(firstOctet, data...)
		data = append(make([]byte, 8-width), data...)
	}

	// The VINT_DATA value MUST be expressed as a big-endian unsigned integer.
	return binary.BigEndian.Uint64(data), nil
}

func Read(reader io.Reader) (uint64, error) {
	return read(reader, true)
}

func ReadWithMarker(reader io.Reader) (uint64, error) {
	return read(reader, false)
}

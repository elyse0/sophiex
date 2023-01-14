package hls

import "sophiex/internal/fragment"

type Fragment struct {
	Generic       fragment.Generic
	MediaSequence int
	Discontinuity int
	Duration      int64 // Milliseconds
	Start         int64 // Unix milliseconds
	End           int64 // Unix milliseconds
}

type Initialization struct {
	Url       string
	ByteRange fragment.ByteRange
}

func (initialization *Initialization) IsEmpty() bool {
	if initialization.Url == "" {
		return true
	}

	return false
}

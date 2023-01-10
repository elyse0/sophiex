package fragment

type ByteRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

func (byteRange *ByteRange) IsEmpty() bool {
	if byteRange.Start == 0 && byteRange.End == 0 {
		return true
	}

	return false
}

type Fragment interface {
	Url() string
	ByteRange() ByteRange
}

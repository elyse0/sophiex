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

type Decryption struct {
	Method string `json:"method"`
	Uri    string `json:"uri"`
	IV     []byte `json:"iv"`
}

func (decryption *Decryption) IsEmpty() bool {
	if decryption.Uri == "" {
		return true
	}

	return false
}

type Generic struct {
	Url        string
	ByteRange  ByteRange
	Decryption Decryption
}

func (generic *Generic) IsEmpty() bool {
	if generic.Url == "" {
		return true
	}

	return false
}

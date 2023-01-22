package components

// Header https://www.rfc-editor.org/rfc/rfc8794.html#section-8.1
// The EBML Header is a declaration that provides processing instructions and identification of the EBML Body.
type Header struct {
	DocType        string `json:"docType"`
	DocTypeVersion uint64 `json:"docTypeVersion"`
}

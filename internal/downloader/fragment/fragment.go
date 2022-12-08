package fragment

import "net/http"

type FragmentRequest struct {
	Index int
	Url   string
}

type FragmentResponse struct {
	Index    int
	Response *http.Response
}

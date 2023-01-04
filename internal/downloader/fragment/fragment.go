package fragment

import "sophiex/internal/parser"

type FragmentRequest struct {
	Index    int
	Fragment parser.HlsFragment
}

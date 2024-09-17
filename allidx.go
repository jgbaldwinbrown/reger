package reger

import (
	"io"
	"regexp"
)

// Returns a [][]int containing indices for all matches in the reader or nil if
// there are no matches. Unlike the equivalent method of *reger.Reger, this does not
// allocate a buffer of bytes.
func FindRuneReaderAllIndex(re *regexp.Regexp, r io.RuneReader) [][]int {
	curlen := 0
	var out [][]int
	for {
		prevlen := curlen
		idxs := re.FindReaderIndex(r)
		if idxs == nil {
			break
		}
		curlen += idxs[1]
		idxs[0] += prevlen
		idxs[1] += prevlen
		out = append(out, idxs)
	}
	return out
}

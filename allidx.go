package reger

import (
	"io"
	"regexp"
	"bufio"
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

// Returns a [][]int containing indices for all submatches for each match in the reader or nil if
// there are no matches. Unlike the equivalent method of *reger.Reger, this does not
// allocate a buffer of bytes.
func FindRuneReaderAllSubmatchIndex(re *regexp.Regexp, r io.RuneReader) [][]int {
	curlen := 0
	var out [][]int
	for {
		prevlen := curlen
		idxs := re.FindReaderSubmatchIndex(r)
		if idxs == nil {
			break
		}
		curlen += idxs[1]
		for i, _ := range idxs {
			idxs[i] += prevlen
		}
		out = append(out, idxs)
	}
	return out
}

func FindReaderAllIndex(re *regexp.Regexp, r io.Reader) [][]int {
	return FindRuneReaderAllIndex(re, bufio.NewReader(r))
}

func FindReaderAllSubmatchIndex(re *regexp.Regexp, r io.Reader) [][]int {
	return FindRuneReaderAllSubmatchIndex(re, bufio.NewReader(r))
}

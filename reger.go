// Reger allows you to use any of the typical regexp.Regexp methods on an
// io.RuneReader or io.Reader. The regexp package only exposes variations on
// FindReaderIndex for io.RuneReaders. The single new type here, Reger,
// maintains an internal buffer of all bytes read via reger.Reger.ReadRune,
// resetting whenever a new regex method is called. This backing buffer is used
// to build all of the byte slices and strings expected by other regex methods.
package reger

import (
	"bufio"
	"io"
	"regexp"
	"unicode/utf8"
)

// The main type here. Holds an io.RuneReader and passes calls to
// reger.Reger.ReadRune through to it, capturing all read runes in an internal buffer.
type Reger struct {
	r   io.RuneReader
	buf []byte
}

// Reads one rune from Reger's captured io.RuneReader, stores it in the
// internal buffer, and returns it.
func (x *Reger) ReadRune() (r rune, size int, err error) {
	r, size, err = x.r.ReadRune()
	if err != nil {
		return 0, 0, nil
	}
	x.buf = utf8.AppendRune(x.buf, r)
	return r, size, nil
}

// Return the bytes read since the last call to a regex method or r.ResetBuffer.
func (r *Reger) Bytes() []byte {
	return r.buf
}

// Empties the internal buffer.
func (r *Reger) ResetBuffer() {
	r.buf = r.buf[:0]
}

// Creates a new Reger from an io.RuneReader. No buffering is done, so no extra
// bytes are read when finding a regex. The reader can still be used by other
// functions without strange outcomes.
func NewReger(r io.RuneReader) *Reger {
	return &Reger{r: r}
}

// Creates a new Reger from an io.Reader. Internally, a *bufio.Reader is
// created, and this buffers the input reader. Do not use this reader with
// other functions.
func NewReaderReger(r io.Reader) *Reger {
	return &Reger{r: bufio.NewReader(r)}
}

// Return a byte slice containing the match, or nil if there is no match.  The
// backing array for this slice will be overwritten when the next regex method is called and
// should be copied if you want to use it later.
func (r *Reger) FindReader(re *regexp.Regexp) []byte {
	r.buf = r.buf[:0]
	idxs := re.FindReaderIndex(r)
	if idxs == nil {
		return nil
	}
	return r.buf[idxs[0]:idxs[1]]
}

// Return a string containing a match, or "" if there is none.
func (r *Reger) FindReaderString(re *regexp.Regexp) string {
	b := r.FindReader(re)
	if b == nil {
		return ""
	}
	return string(b)
}

// Returns a [][]byte containing all submatches, or nil if there is no match.
// The backing array for these will be overwritten when the next regex method is called and
// should be copied if you want to use it later.
func (r *Reger) FindReaderSubmatch(re *regexp.Regexp) [][]byte {
	r.buf = r.buf[:0]
	idxs := re.FindReaderSubmatchIndex(r)
	if idxs == nil {
		return nil
	}
	out := make([][]byte, 0, len(idxs)/2)
	for i := 0; i < len(idxs); i += 2 {
		idxpair := idxs[i : i+2]
		out = append(out, r.buf[idxpair[0]:idxpair[1]])
	}
	return out
}

// Returns a []string containing all submatches or nil if there is no match
func (r *Reger) FindReaderStringSubmatch(re *regexp.Regexp) []string {
	bs := r.FindReaderSubmatch(re)
	if bs == nil {
		return nil
	}
	out := make([]string, 0, len(bs))
	for _, b := range bs {
		if b == nil {
			out = append(out, "")
		} else {
			out = append(out, string(b))
		}
	}
	return out
}

// Returns a [][]int containing indices for all matches in the reader or nil if
// there are no matches.
func (r *Reger) FindReaderAllIndex(re *regexp.Regexp) [][]int {
	r.buf = r.buf[:0]
	var out [][]int
	for {
		prevlen := len(r.buf)
		idxs := re.FindReaderIndex(r)
		if idxs == nil {
			break
		}
		idxs[0] += prevlen
		idxs[1] += prevlen
		out = append(out, idxs)
	}
	return out
}

// Returns a [][]byte containing all matches in the reader or nil if there are
// no matches.  The backing array for these will be overwritten when the next
// regex method is called and should be copied if you want to use it later.
func (r *Reger) FindReaderAll(re *regexp.Regexp) [][]byte {
	idxs := r.FindReaderAllIndex(re)
	if idxs == nil {
		return nil
	}
	out := make([][]byte, 0, len(idxs))
	for _, idxpair := range idxs {
		if idxpair == nil {
			out = append(out, nil)
		} else {
			out = append(out, r.buf[idxpair[0]:idxpair[1]])
		}
	}
	return out
}

// Returns a []string containing all matches in the reader or nil if there are
// no matches.
func (r *Reger) FindReaderAllString(re *regexp.Regexp) []string {
	var out []string
	for ss := r.FindReaderStringSubmatch(re); ss != nil; ss = r.FindReaderStringSubmatch(re) {
		out = append(out, ss[0])
	}
	return out
}

// Returns a [][]int containing submatch indices for all matches in the reader
// or nil if there are no matches.
func (r *Reger) FindReaderAllSubmatchIndex(re *regexp.Regexp) [][]int {
	r.buf = r.buf[:0]
	var out [][]int
	for {
		prevlen := len(r.buf)
		idxs := re.FindReaderSubmatchIndex(r)
		if idxs == nil {
			break
		}
		for i, _ := range idxs {
			idxs[i] += prevlen
		}
		out = append(out, idxs)
	}
	return out
}

// Returns a [][][]byte containing submatches for all matches in the reader or
// nil if there are no matches.  The backing array for these will be
// overwritten when the next regex method is called and should be copied if you
// want to use it later.
func (r *Reger) FindReaderAllSubmatch(re *regexp.Regexp) [][][]byte {
	idxs := r.FindReaderAllSubmatchIndex(re)
	if idxs == nil {
		return nil
	}
	out := make([][][]byte, 0, len(idxs))
	for _, idxset := range idxs {
		if idxset == nil {
			out = append(out, nil)
			continue
		}

		bset := make([][]byte, 0, len(idxset)/2)
		for i := 0; i < len(idxset); i += 2 {
			idxpair := idxset[i : i+2]
			bset = append(bset, r.buf[idxpair[0]:idxpair[1]])
		}
		out = append(out, bset)
	}
	return out
}

// Returns a [][]string containing submatches for all matches in the reader or
// nil if there are no matches.
func (r *Reger) FindReaderAllStringSubmatch(re *regexp.Regexp) [][]string {
	var out [][]string
	for ss := r.FindReaderStringSubmatch(re); ss != nil; ss = r.FindReaderStringSubmatch(re) {
		out = append(out, ss)
	}
	return out
}

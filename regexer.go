package regexer

import (
	"regexp"
	"io"
	"unicode/utf8"
	"bufio"
)

type Regexer struct {
	r io.RuneReader
	buf []byte
}

func (x *Regexer) ReadRune() (r rune, size int, err error) {
	r, size, err = x.r.ReadRune()
	if err != nil {
		return 0, 0, nil
	}
	x.buf = utf8.AppendRune(x.buf, r)
	return r, size, nil
}

func NewRegexer(r io.RuneReader) *Regexer {
	return &Regexer{r: r}
}

func NewReaderRegexer(r io.Reader) *Regexer {
	return &Regexer{r: bufio.NewReader(r)}
}

func (r *Regexer) FindReader(re *regexp.Regexp) []byte {
	r.buf = r.buf[:0]
	idxs := re.FindReaderIndex(r)
	if idxs == nil {
		return nil
	}
	return r.buf[idxs[0]:idxs[1]]
}

func (r *Regexer) FindReaderString(re *regexp.Regexp) string {
	b := r.FindReader(re)
	if b == nil {
		return ""
	}
	return string(b)
}

func (r *Regexer) FindReaderSubmatch(re *regexp.Regexp) [][]byte {
	r.buf = r.buf[:0]
	idxs := re.FindReaderSubmatchIndex(r)
	if idxs == nil {
		return nil
	}
	out := make([][]byte, 0, len(idxs) / 2)
	for i := 0; i < len(idxs); i += 2 {
		idxpair := idxs[i:i + 2]
		out = append(out, r.buf[idxpair[0]:idxpair[1]])
	}
	return out
}

func (r *Regexer) FindReaderStringSubmatch(re *regexp.Regexp) []string {
	bs := r.FindReaderSubmatch(re)
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

func (r *Regexer) FindReaderAllIndex(re *regexp.Regexp) [][]int {
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

func (r *Regexer) FindReaderAll(re *regexp.Regexp) [][]byte {
	idxs := r.FindReaderAllIndex(re)
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

func (r *Regexer) FindReaderAllString(re *regexp.Regexp) []string {
	bs := r.FindReaderAll(re)
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

func (r *Regexer) FindReaderAllSubmatchIndex(re *regexp.Regexp) [][]int {
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

func (r *Regexer) FindReaderAllSubmatch(re *regexp.Regexp) [][][]byte {
	idxs := r.FindReaderAllSubmatchIndex(re)
	out := make([][][]byte, 0, len(idxs))
	for _, idxset := range idxs {
		if idxset == nil {
			out = append(out, nil)
			continue
		}

		bset := make([][]byte, 0, len(idxset) / 2)
		for i := 0; i < len(idxset); i += 2 {
			idxpair := idxset[i:i + 2]
			bset = append(bset, r.buf[idxpair[0]:idxpair[1]])
		}
		out = append(out, bset)
	}
	return out
}

func (r *Regexer) FindReaderAllStringSubmatch(re *regexp.Regexp) [][]string {
	bsets := r.FindReaderAllSubmatch(re)
	if bsets == nil {
		return nil
	}
	out := make([][]string, 0, len(bsets))
	for _, bset := range bsets {
		outset := make([]string, 0, len(bset))
		for _, b := range bset {
			outset = append(outset, string(b))
		}
		out = append(out, outset)
	}
	return out
}

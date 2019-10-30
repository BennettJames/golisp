package golisp2

import (
	"bufio"
	"io"
)

// RuneScanner is used to iteratively read a source for characters.
type RuneScanner struct {
	err error
	r   rune
	buf *bufio.Reader
}

// NewRuneScanner initializes a RuneScanner around the given string.
func NewRuneScanner(src io.Reader) *RuneScanner {
	return &RuneScanner{
		buf: bufio.NewReader(src),
	}
}

// Rune returns the rune at the current index in the scanner.
func (rs *RuneScanner) Rune() rune {
	return rs.r
}

// Advance moves the scanner ahead one value,
func (rs *RuneScanner) Advance() {
	// note (bs): technically possible this value is not valid and a
	// unicode.ReplacementChar is returned. If so, possible that should be handled
	// here.
	if rs.err != nil {
		return
	}
	r, _, err := rs.buf.ReadRune()
	if err != nil {
		rs.err = err
		return
	}
	rs.r = r
}

// Done indicates if the scanner has reached completion.
func (rs *RuneScanner) Done() bool {
	return rs.err != nil
}

func (rs *RuneScanner) Err() error {
	return rs.err
}

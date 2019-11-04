package golisp2

import (
	"bufio"
	"io"
)

type (
	// RuneScanner is used to iteratively read a source for characters.
	RuneScanner struct {
		err error
		r   rune
		pos ScannerPosition
		buf *bufio.Reader
	}

	// ScannerPosition contains location information for runes and tokens.
	ScannerPosition struct {
		SourceFile string
		Col, Row   int
	}
)

// NewRuneScanner initializes a RuneScanner around the given string.
func NewRuneScanner(srcName string, src io.Reader) *RuneScanner {
	return &RuneScanner{
		buf: bufio.NewReader(src),
		pos: ScannerPosition{
			SourceFile: srcName,
			Row:        1,
		},
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
	//
	// Also: I'm sorta inclined to deliberately error out on \0. It's really weird
	// for that to be in a source file, and I sorta implicitly use '\0' here as a
	// zero value to indicate either done or uninitialized.
	if rs.err != nil {
		return
	}
	r, _, err := rs.buf.ReadRune()
	if err != nil {
		rs.err = err
		rs.r = 0
		return
	}
	if rs.r == '\n' {
		rs.pos.Row++
		rs.pos.Col = 1
	} else {
		rs.pos.Col++
	}
	rs.r = r
}

// Pos returns the current location of the scanner relative to it's source.
func (rs *RuneScanner) Pos() ScannerPosition {
	return rs.pos
}

// Done indicates if the scanner has reached completion.
func (rs *RuneScanner) Done() bool {
	return rs.err != nil
}

// Err returns any error encountered during the scan. Will be io.EOF if the
// scanner simply completed.
func (rs *RuneScanner) Err() error {
	return rs.err
}

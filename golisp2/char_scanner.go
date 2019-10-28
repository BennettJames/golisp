package golisp2

// CharScanner is used to iteratively read a source for characters.
type CharScanner struct {
	index int
	runes []rune
}

// NewCharScanner initializes a CharScanner around the given string.
func NewCharScanner(str string) *CharScanner {
	runes := []rune{}
	for _, r := range str {
		runes = append(runes, r)
	}

	return &CharScanner{
		runes: runes,
	}
}

// Next returns the next rune in the scanner, and advances progress.
func (cs *CharScanner) Next() (rune, bool) {
	r, ok := cs.Peek()
	if ok {
		cs.index++
	}
	return r, ok
}

// Peek returns the next rune, if available. Returns false if not.
func (cs *CharScanner) Peek() (rune, bool) {
	if cs.Done() {
		return 0, false
	}
	r := cs.runes[cs.index]
	return r, true
}

// Done indicates if the scanner has reached completion.
func (cs *CharScanner) Done() bool {
	return cs.index >= len(cs.runes)
}

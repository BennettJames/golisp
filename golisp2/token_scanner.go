package golisp2

import (
	"unicode"
)

type (
	// TokenScanner reads over an input source of characters, transforming them
	// into tokens.
	TokenScanner struct {
		done bool
		t    *ScannedToken
		st   *subTokenScanner
	}

	// subTokenScanner is a private substructure for TokenScanner that does most
	// of the work. It's responsible for buffering in-progress tokens.
	subTokenScanner struct {
		src      *RuneScanner
		buf      []byte
		startPos ScannerPosition
	}
)

// NewTokenScanner creates a new TokenScanner around the provided source.
func NewTokenScanner(src *RuneScanner) *TokenScanner {
	return &TokenScanner{
		st: newSubTokenScanner(src),
	}
}

// Done indicates if the underlying source has been exhausted, with no more
// values to read.
func (ts *TokenScanner) Done() bool {
	return ts.done
}

// Err returns any error encountered while scanning the input. Will be io.EOF if
// the scan completed the input.
func (ts *TokenScanner) Err() error {
	return ts.st.src.Err()
}

// Advance will read in the next token into the scanner.
func (ts *TokenScanner) Advance() {
	var maybeNextT *ScannedToken
	for !ts.st.src.Done() {
		maybeNextT = scanNextToken(ts.st)
		if maybeNextT != nil && maybeNextT.Typ == CommentTT {
			// skip comments; by definition they don't need to be parsed
			continue
		}
		break
	}
	ts.t = maybeNextT
	if maybeNextT == nil {
		ts.done = true
	}
}

// Token returns the token currently read by the scanner. Will be nil if
// `Advance` has never been called.
func (ts *TokenScanner) Token() *ScannedToken {
	return ts.t
}

func newSubTokenScanner(src *RuneScanner) *subTokenScanner {
	return &subTokenScanner{
		src: src,
	}
}

func (ss *subTokenScanner) Done() bool {
	// note (bs): not *quite* sure yet if this is correct, but o.k. for now.
	// Particularly, what if the underlying stream is done but there's still a
	// token prepped for grabbing? Need to make sure that I neither double-process
	// or skip the last rune.
	return ss.src.Done()
}

// Rune returns the current rune being scanned.
func (ss *subTokenScanner) Rune() rune {
	return ss.src.Rune()
}

// Advance adds the current rune to the buffer, and moves to the next step.
func (ss *subTokenScanner) Advance() {
	if ss.src.Done() {
		return
	}
	if len(ss.buf) == 0 {
		ss.startPos = ss.src.Pos()
	}
	ss.buf = append(ss.buf, []byte(string(ss.src.Rune()))...)
	ss.src.Advance()
}

// Skip will advance to the next rune, but without including it in the buffer.
func (ss *subTokenScanner) Skip() {
	if ss.src.Done() {
		return
	}
	ss.src.Advance()
}

// Complete drains the buffer and returns a token consisting of the type and the
// buffer as the value.
func (ss *subTokenScanner) Complete(t TokenType) *ScannedToken {
	val := string(ss.buf)
	ss.buf = nil
	return &ScannedToken{
		Typ:   t,
		Value: val,
		Pos:   ss.startPos,
	}
}

// FlushInvalid writes the current rune to the buffer, and completes the scan
// with an invlaid type. Useful for cases where the current rune is unscannable;
// and the only thing to do is to advance and mark it invalid.
func (ss *subTokenScanner) FlushInvalid() *ScannedToken {
	ss.Advance()
	return ss.Complete(InvalidTT)
}

func scanNextToken(s *subTokenScanner) *ScannedToken {
	// Remove any leading whitespace
	for !s.src.Done() && (s.src.Rune() == '\x00' || unicode.IsSpace(s.src.Rune())) {
		s.Skip()
	}
	if s.src.Done() {
		return nil
	}

	if s.Rune() == '(' {
		s.Advance()
		return s.Complete(OpenParenTT)
	} else if s.Rune() == ')' {
		s.Advance()
		return s.Complete(CloseParenTT)
	} else if s.Rune() == ';' {
		return tryLexComment(s)
	} else if s.Rune() == '-' {
		return tryLexSignedValue(s)
	} else if isOperatorRune(s.Rune()) {
		return tryLexOperator(s)
	} else if isDigitRune(s.Rune()) {
		return tryLexNumber(s)
	} else if isDoubleQuoteRune(s.Rune()) {
		return tryLexString(s)
	} else if isIdentStartRune(s.Rune()) {
		return tryLexIdent(s)
	}

	return s.FlushInvalid()
}

func tryLexComment(s *subTokenScanner) *ScannedToken {
	if s.Rune() != ';' {
		return s.FlushInvalid()
	}
	s.Advance()
	for !s.Done() && s.Rune() != '\n' {
		s.Advance()
	}
	return s.Complete(CommentTT)
}

func tryLexSignedValue(s *subTokenScanner) *ScannedToken {
	if s.Rune() != '-' {
		return s.FlushInvalid()
	}
	s.Advance()
	if isDigitRune(s.Rune()) {
		return tryLexNumber(s)
	}
	return tryLexOperatorTail(s)
}

func tryLexOperator(s *subTokenScanner) *ScannedToken {
	if !isOperatorRune(s.Rune()) {
		return s.FlushInvalid()
	}
	s.Advance()
	return tryLexOperatorTail(s)
}

func tryLexOperatorTail(s *subTokenScanner) *ScannedToken {
	for {
		if isOperatorRune(s.Rune()) {
			s.Advance()
			continue
		}
		if scannerAtBoundary(s) {
			return s.Complete(OpTT)
		}
		return s.FlushInvalid()
	}
}

func tryLexNumber(s *subTokenScanner) *ScannedToken {
	// note (bs): this is a more general problem; but I think ensuring
	// "at-least-one-digit" like this is pretty clumsy. Maybe there should be a
	// generic way to "slurp down" chars of least a given length.
	if !unicode.IsDigit(s.Rune()) {
		return s.FlushInvalid()
	}
	s.Advance()

	// note (bs): this still isn't *great* as far as division of responsibilities
	// is concerned. May want a somewhat easier way to do things like specify a
	// minimum number of digits to lex in a pass.
	for {
		if isDigitRune(s.Rune()) {
			s.Advance()
			continue
		}

		// note (bs): this isn't technically correct, as it could tolerate multiple
		// decimal points. Need to subdivide further for this to be right.
		if isDecimalRune(s.Rune()) {
			s.Advance()
			if isDigitRune(s.Rune()) {
				s.Advance()
				continue
			}
			return s.Complete(InvalidTT)
		}

		if scannerAtBoundary(s) {
			return s.Complete(NumberTT)
		}
		return s.FlushInvalid()
	}
}

func tryLexString(s *subTokenScanner) *ScannedToken {
	if !isDoubleQuoteRune(s.Rune()) {
		return s.FlushInvalid()
	}
	s.Advance()

	for {
		if s.Done() || isNewlineRune(s.Rune()) {
			return s.FlushInvalid()
		}

		if isDoubleQuoteRune(s.Rune()) {
			s.Advance()
			if scannerAtBoundary(s) {
				return s.Complete(StringTT)
			}
			return s.FlushInvalid()
		}

		// todo (bs): need to process escaped characters here. That will require a
		// bit of an API shift: the tokenization process will need to be a little
		// more state-dependent and you'd need to be able to push specific runes
		// onto the scanner; rather than strictly work off of a process of
		// "advance".
		//
		// I suspect that this whole process should be replaced by some better
		// sub-tokenization for each top-level element. I'd guess at that point,
		// finer-grained sub-parsing like that would be more appropriate. Right now,
		// the tokenization is very crude; and doesn't distinguish quotes from the
		// rest of the inner elements.

		s.Advance()
	}
}

func tryLexIdent(s *subTokenScanner) *ScannedToken {
	if !isIdentStartRune(s.Rune()) {
		return s.FlushInvalid()
	}
	s.Advance()

	for {
		if scannerAtBoundary(s) {
			return s.Complete(IdentTT)
		}
		if isIdentRune(s.Rune()) {
			s.Advance()
			continue
		}
		return s.FlushInvalid()
	}
}

func scannerAtBoundary(s *subTokenScanner) bool {
	return s.Done() ||
		isSpaceRune(s.Rune()) ||
		isOpenParenRune(s.Rune()) ||
		isCloseParenRune(s.Rune())
}

func isDigitRune(r rune) bool {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	default:
		return false
	}
}

func isOperatorRune(r rune) bool {
	switch r {
	case '-', '+', '/', '*', '&', '^', '%', '!', '|', '<', '>', '=':
		return true
	default:
		return false
	}
}

func isDecimalRune(r rune) bool {
	return r == '.'
}

func isSpaceRune(r rune) bool {
	return unicode.IsSpace(r)
}

func isOpenParenRune(r rune) bool {
	return r == '('
}

func isCloseParenRune(r rune) bool {
	return r == ')'
}

func isDoubleQuoteRune(r rune) bool {
	return r == '"'
}

func isNewlineRune(r rune) bool {
	return r == '\n'
}

func isIdentStartRune(r rune) bool {
	// note (bs): not sure if this adequate/correct; but will probably be o.k. to
	// start
	return unicode.IsLetter(r)
}

func isIdentRune(r rune) bool {
	return isIdentStartRune(r) || unicode.IsDigit(r)
}

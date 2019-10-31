package golisp2

import (
	"strings"
	"unicode"
)

// TokenizeString converts the provided string to a list of tokens.
func TokenizeString(str string) []ScannedToken {
	tokens := []ScannedToken{}

	cs := NewRuneScanner(strings.NewReader(str))
	ts := NewTokenScanner(cs)
	for !ts.Done() {
		nextT := ts.Next()
		if nextT == nil {
			break
		}
		tokens = append(tokens, *nextT)
		if nextT.Typ == InvalidTT {
			break
		}
	}
	return tokens
}

type (
	// TokenScanner reads over an input source of characters, transforming them
	// into tokens.
	TokenScanner struct {
		st *subTokenScanner
	}

	// subTokenScanner is a private substructure for TokenScanner that does most
	// of the work. It's responsible for buffering in-progress tokens.
	subTokenScanner struct {
		src *RuneScanner
		buf []byte
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
	return ts.st.Done()
}

// Err returns any error encountered while scanning the input. Will be io.EOF if
// the scan completed the input.
func (ts *TokenScanner) Err() error {
	return ts.st.src.Err()
}

// Next will read in from the source until a new token is discovered. If the
// source is empty, returns nil.
func (ts *TokenScanner) Next() *ScannedToken {
	return scanNextToken(ts.st)
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
	}
}

func scanNextToken(s *subTokenScanner) *ScannedToken {
	// Removed any leading whitespace
	//
	// note (bs): not sure the handling of null bytes here is really correct
	for !s.src.Done() && (s.src.Rune() == '\x00' || unicode.IsSpace(s.src.Rune())) {
		s.src.Advance()
	}
	if s.src.Done() {
		return nil
	}

	// note (bs): this is kinda inefficient - if you hoisted the initial chars,
	// you could always map to the next function. Oh well.
	tryParsers := []func(*subTokenScanner) *ScannedToken{
		tryParseOpenParen,
		tryParseCloseParen,
		tryParseSignedValue,
		tryParseOperator,
		tryParseNumber,
		tryParseString,
		tryParseIdent,
	}

	for _, p := range tryParsers {
		if t := p(s); t != nil {
			return t
		}
	}

	return s.Complete(InvalidTT)
}

func tryParseOpenParen(s *subTokenScanner) *ScannedToken {
	if s.Rune() == '(' {
		s.Advance()
		return s.Complete(OpenParenTT)
	}
	return nil
}

func tryParseCloseParen(s *subTokenScanner) *ScannedToken {
	if s.Rune() == ')' {
		s.Advance()
		return s.Complete(CloseParenTT)
	}
	return nil
}

func tryParseSignedValue(s *subTokenScanner) *ScannedToken {
	if s.Rune() != '-' {
		return nil
	}
	s.Advance()
	if scannerAtDigit(s) {
		s.Advance()
		return parseNumber(s)
	}
	return parseOperator(s)
}

func tryParseOperator(s *subTokenScanner) *ScannedToken {
	if scannerAtOperator(s) {
		s.Advance()
		return parseOperator(s)
	}
	return nil
}

func parseOperator(s *subTokenScanner) *ScannedToken {
	for {
		if scannerAtOperator(s) {
			s.Advance()
			continue
		}
		if scannerAtBoundary(s) {
			return s.Complete(OpTT)
		}
		s.Advance()
		return s.Complete(InvalidTT)
	}
}

func tryParseNumber(s *subTokenScanner) *ScannedToken {
	// note (bs): just does nonnegative integers right now; not very
	// sophisticated.
	if scannerAtDigit(s) {
		s.Advance()
		return parseNumber(s)
	}
	return nil
}

func parseNumber(s *subTokenScanner) *ScannedToken {
	// note (bs): this still isn't *great* as far as division of responsibilities
	// is concerned. May want a somewhat easier way to do things like specify a
	// minimum number of digits to parse in a pass.
	for {
		if scannerAtDigit(s) {
			s.Advance()
			continue
		}

		// note (bs): this isn't technically correct, as it could tolerate multiple
		// decimal points. Need to subdivide further for this to be right.
		if scannerAtDecimal(s) {
			s.Advance()
			if scannerAtDigit(s) {
				s.Advance()
				continue
			}
			return s.Complete(InvalidTT)
		}

		if scannerAtBoundary(s) {
			return s.Complete(NumberTT)
		}
		return s.Complete(InvalidTT)
	}
}

func tryParseString(s *subTokenScanner) *ScannedToken {
	if scannerAtDoubleQuote(s) {
		s.Advance()
		return parseString(s)
	}
	return nil
}

func parseString(s *subTokenScanner) *ScannedToken {
	for {
		if s.Done() || scannerAtNewline(s) {
			return s.Complete(InvalidTT)
		}

		if scannerAtDoubleQuote(s) {
			s.Advance()
			if scannerAtBoundary(s) {
				return s.Complete(StringTT)
			}
			return s.Complete(InvalidTT)
		}

		// todo (bs): need to process escaped characters here. That will require a
		// bit of an API shift: the tokenization process will need to
		//
		// I suspect that this whole process should be replaced by some better
		// sub-tokenization for each top-level element. I'd guess at that point,
		// finer-grained sub-parsing like that would be more appropriate. Right now,
		// the tokenization is very crude; and doesn't distinguish quotes from the
		// rest of the inner elements.

		s.Advance()
	}
}

func tryParseIdent(s *subTokenScanner) *ScannedToken {
	if scannerAtIdentStart(s) {
		s.Advance()
		return parseIdent(s)
	}
	return nil
}

func parseIdent(s *subTokenScanner) *ScannedToken {
	for {
		if scannerAtBoundary(s) {
			return s.Complete(IdentTT)
		}
		if scannerAtIdent(s) {
			s.Advance()
			continue
		}
		s.Advance()
		return s.Complete(InvalidTT)
	}
}

func scannerAtBoundary(s *subTokenScanner) bool {
	return s.Done() ||
		scannerAtSpace(s) ||
		scannerAtOpenParen(s) ||
		scannerAtCloseParen(s)
}

func scannerAtDigit(s *subTokenScanner) bool {
	if s.Done() {
		return false
	}
	switch s.Rune() {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	default:
		return false
	}
}

func scannerAtOperator(s *subTokenScanner) bool {
	if s.Done() {
		return false
	}
	switch s.Rune() {
	case '-', '+', '/', '*', '&', '^', '%', '!', '|', '<', '>', '=':
		return true
	default:
		return false
	}
}

func scannerAtDecimal(s *subTokenScanner) bool {
	if s.Done() {
		return false
	}
	return s.Rune() == '.'
}

func scannerAtSpace(s *subTokenScanner) bool {
	if s.Done() {
		return false
	}
	return unicode.IsSpace(s.Rune())
}

func scannerAtOpenParen(s *subTokenScanner) bool {
	if s.Done() {
		return false
	}
	return s.Rune() == '('
}

func scannerAtCloseParen(s *subTokenScanner) bool {
	if s.Done() {
		return false
	}
	return s.Rune() == ')'
}

func scannerAtDoubleQuote(s *subTokenScanner) bool {
	if s.Done() {
		return false
	}
	return s.Rune() == '"'
}

func scannerAtNewline(s *subTokenScanner) bool {
	if s.Done() {
		return false
	}
	return s.Rune() == '\n'
}

func scannerAtIdentStart(s *subTokenScanner) bool {
	if s.Done() {
		return false
	}
	// note (bs): not sure if this adequate/correct; but will probably be o.k. to
	// start
	return unicode.IsLetter(s.Rune())
}

func scannerAtIdent(s *subTokenScanner) bool {
	if s.Done() {
		return false
	}
	return scannerAtIdentStart(s) ||
		unicode.IsDigit(s.Rune())
}

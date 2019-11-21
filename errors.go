package golisp2

import "fmt"

type (
	// ParseError reflects an error that took place during parsing. It contains
	// information
	ParseError struct {
		msg   string
		token ScannedToken
	}

	// ForbiddenRuneError indicates that an illegal character was found in the
	// source.
	ForbiddenRuneError struct {
		r   rune
		pos ScannerPosition
	}

	// TypeError is a runtime error when the incorrect type is passed to a
	// function.
	TypeError struct {
		actual, expected string
		pos              ScannerPosition
	}
)

// NewParseError creates a new parse error with the given message and token.
func NewParseError(msg string, token ScannedToken) *ParseError {
	return &ParseError{
		msg:   msg,
		token: token,
	}
}

// NewParseEOFError represents a parsing error for unexpected EOF.
func NewParseEOFError(msg string, pos ScannerPosition) *ParseError {
	return &ParseError{
		msg: msg,
		token: ScannedToken{
			Typ: NoTT,
			Pos: pos,
		},
	}
}

// Error returns the informational error string about the parse error.
func (pe ParseError) Error() string {
	// note (bs): I don't think this is a very well-laid out error message, but
	// it's a place to start at least.
	msg, token, pos := pe.msg, pe.token, pe.token.Pos
	return fmt.Sprintf(
		"Parse error %s for token `%s`: file '%s' at line %d, column %d",
		msg, token.Value, pos.SourceFile, pos.Row, pos.Col)
}

// NewForbiddenRuneError creates a ForbiddenRuneError for the given rune and
// location it was found.
func NewForbiddenRuneError(r rune, pos ScannerPosition) *ForbiddenRuneError {
	return &ForbiddenRuneError{
		r:   r,
		pos: pos,
	}
}

// Error returns the informational error string about the parse error.
func (pe ForbiddenRuneError) Error() string {
	return fmt.Sprintf(
		"Forbidden rune '%x' found in scan of '%s' (line %d, col %d)",
		pe.r, pe.pos.SourceFile, pe.pos.Row, pe.pos.Col)
}

// NewTypeError creates a new type error with the actual and expected types at
// the given location in source.
func NewTypeError(actual, expected string, pos ScannerPosition) *TypeError {
	return &TypeError{
		actual:   actual,
		expected: expected,
		pos:      pos,
	}
}

func (te TypeError) Error() string {
	return fmt.Sprintf(
		"Type error: expected '%s', got '%s' (%s:%d)",
		te.expected, te.actual,
		te.pos.SourceFile, te.pos.Row)
}

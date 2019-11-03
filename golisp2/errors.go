package golisp2

import "fmt"

type (
	// ParseError reflects an error that took place during parsing. It contains
	// information
	ParseError struct {
		msg   string
		token ScannedToken
	}
)

// NewParseError creates a new parse error with the given message and token.
func NewParseError(msg string, token ScannedToken) *ParseError {
	return &ParseError{
		msg:   msg,
		token: token,
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

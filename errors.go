package golisp2

import "fmt"

type (
	// ParseError reflects an error that took place during parsing. It contains
	// information
	ParseError struct {
		Msg   string
		Token ScannedToken
	}

	// ForbiddenRuneError indicates that an illegal character was found in the
	// source.
	ForbiddenRuneError struct {
		R   rune
		Pos ScannerPosition
	}

	// TypeError is a runtime error when the incorrect type is passed to a
	// function.
	TypeError struct {
		Actual, Expected string
		Pos              ScannerPosition
	}

	// EvalError is a basic runtime error indicating something went wrong during
	// execution.
	EvalError struct {
		Msg string
		Pos ScannerPosition
	}

	// ArgTypeError indicates a mismatch between a given argument value and the
	// expected type.
	//
	// note (bs): this is a fairly awkward error type that probably should just be
	// a type error. However, there are some structural limitations with built-ins
	// that make that challenging.
	ArgTypeError struct {
		FnName           string
		ArgI             int
		Expected, Actual string
	}
)

// NewParseError creates a new parse error with the given message and token.
func NewParseError(msg string, token ScannedToken) *ParseError {
	return &ParseError{
		Msg:   msg,
		Token: token,
	}
}

// NewParseEOFError represents a parsing error for unexpected EOF.
func NewParseEOFError(msg string, pos ScannerPosition) *ParseError {
	return &ParseError{
		Msg: msg,
		Token: ScannedToken{
			Typ: NoTT,
			Pos: pos,
		},
	}
}

// Error returns the informational error string about the parse error.
func (pe ParseError) Error() string {
	// note (bs): I don't think this is a very well-laid out error message, but
	// it's a place to start at least.
	msg, token, pos := pe.Msg, pe.Token, pe.Token.Pos
	return fmt.Sprintf(
		"Parse error %s for token `%s`: file '%s' at line %d, column %d",
		msg, token.Value, pos.SourceFile, pos.Row, pos.Col)
}

// NewForbiddenRuneError creates a ForbiddenRuneError for the given rune and
// location it was found.
func NewForbiddenRuneError(r rune, pos ScannerPosition) *ForbiddenRuneError {
	return &ForbiddenRuneError{
		R:   r,
		Pos: pos,
	}
}

// Error returns the informational error string about the parse error.
func (pe ForbiddenRuneError) Error() string {
	return fmt.Sprintf(
		"Forbidden rune '%x' found in scan of '%s' (line %d, col %d)",
		pe.R, pe.Pos.SourceFile, pe.Pos.Row, pe.Pos.Col)
}

// NewTypeError creates a new type error with the actual and expected types at
// the given location in source.
func NewTypeError(actual, expected string, pos ScannerPosition) *TypeError {
	return &TypeError{
		Actual:   actual,
		Expected: expected,
		Pos:      pos,
	}
}

func (te TypeError) Error() string {
	return fmt.Sprintf(
		"Type error: expected '%s', got '%s' (%s:%d)",
		te.Expected, te.Actual,
		te.Pos.SourceFile, te.Pos.Row)
}

func (ee EvalError) Error() string {
	return fmt.Sprintf("Eval error '%s': '%s' (line %d, col %d)",
		ee.Msg, ee.Pos.SourceFile, ee.Pos.Row, ee.Pos.Col)
}

func (ate *ArgTypeError) Error() string {
	return fmt.Sprintf("Arg-type error in '%s' at arg %d: expected '%s', got '%s'",
		ate.FnName, ate.ArgI, ate.Expected, ate.Actual)
}

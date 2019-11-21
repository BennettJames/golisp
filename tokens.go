package golisp2

import "fmt"

type (
	// TokenType is a basic enum-style type used to denote all the different types
	// of tokens.
	TokenType int

	// ScannedToken contains a pair between token type, and the value scanned by
	// the token.
	ScannedToken struct {
		Typ   TokenType
		Value string
		Pos   ScannerPosition
	}
)

const (
	// NoTT is an empty token type. Should never be used with an actual token; use
	// "InvalidTT" to reflect a value that does not map to a type.
	NoTT TokenType = iota

	// InvalidTT represents a scanned value that is not a valid token type.
	InvalidTT

	// OpenParenTT is a single open parenthese token type.
	OpenParenTT

	// CloseParenTT is a single closed parenthese token type.
	CloseParenTT

	// IdentTT is an identifier token type.
	IdentTT

	// OpTT is an operator token type.
	OpTT

	// NumberTT is a number token type.
	NumberTT

	// StringTT is a string token type.
	StringTT

	// CommentTT represents a comment.
	CommentTT
)

// String is just a simple mapping to a human readable string for token types.
func (tt TokenType) String() string {
	switch tt {
	case NoTT:
		return "NoTT"
	case InvalidTT:
		return "InvalidTT"
	case OpenParenTT:
		return "OpenParenTT"
	case CloseParenTT:
		return "CloseParenTT"
	case IdentTT:
		return "IdentTT"
	case OpTT:
		return "OpTT"
	case NumberTT:
		return "NumberTT"
	case StringTT:
		return "StringTT"
	case CommentTT:
		return "CommentTT"
	default:
		return fmt.Sprintf("<unknown type %d>", tt)
	}
}

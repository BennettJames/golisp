package golisp2

import "testing"

import "github.com/stretchr/testify/require"

func Test_ForbiddenRuneError(t *testing.T) {
	err := NewForbiddenRuneError('\x00', ScannerPosition{
		SourceFile: "abc.l",
		Col:        3,
		Row:        4,
	})
	require.Contains(t, err.Error(), "Forbidden")
}

func Test_ParseError(t *testing.T) {
	err := NewParseError(
		"test parse error",
		ScannedToken{
			Typ:   IdentTT,
			Value: "abc.....efg",
			Pos: ScannerPosition{
				SourceFile: "abc.l",
				Col:        3,
				Row:        4,
			},
		},
	)
	require.Contains(t, err.Error(), "Parse")
}

func Test_TypeError(t *testing.T) {
	err := NewTypeError(
		"number",
		"string",
		ScannerPosition{
			SourceFile: "abc.l",
			Col:        3,
			Row:        4,
		},
	)
	require.Contains(t, err.Error(), "Type")
}

func Test_EvalError(t *testing.T) {
	err := EvalError{
		Msg: "runtime error",
		Pos: ScannerPosition{
			SourceFile: "abc.l",
			Col:        3,
			Row:        4,
		},
	}
	require.Contains(t, err.Error(), "Eval")
}

func Test_ArgTypeError(t *testing.T) {
	err := ArgTypeError{
		FnName:   "add",
		ArgI:     2,
		Expected: "Number",
		Actual:   "nil",
	}
	require.Contains(t, err.Error(), "Arg")
}

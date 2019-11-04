package golisp2

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Tokenization(t *testing.T) {
	fName := "testFile.l"
	makePos := func(c, r int) ScannerPosition {
		return ScannerPosition{
			SourceFile: fName,
			Col:        c,
			Row:        r,
		}
	}

	testCases := []struct {
		Name     string
		Input    string
		Output   []ScannedToken
		Disabled bool
	}{
		{
			Name:  "openClose",
			Input: `(   )`,
			Output: []ScannedToken{
				ScannedToken{
					Typ:   OpenParenTT,
					Value: "(",
				},
				ScannedToken{
					Typ:   CloseParenTT,
					Value: ")",
				},
			},
		},
		{
			Name:  "operators",
			Input: `+ - / * &^%!|<>= =<<<`,
			Output: []ScannedToken{
				ScannedToken{
					Typ:   OpTT,
					Value: "+",
				},
				ScannedToken{
					Typ:   OpTT,
					Value: "-",
				},
				ScannedToken{
					Typ:   OpTT,
					Value: "/",
				},
				ScannedToken{
					Typ:   OpTT,
					Value: "*",
				},
				ScannedToken{
					Typ:   OpTT,
					Value: "&^%!|<>=",
				},
				ScannedToken{
					Typ:   OpTT,
					Value: "=<<<",
				},
			},
		},
		{
			Name:  "basicNumbers",
			Input: `1 57.123 -2`,
			Output: []ScannedToken{
				ScannedToken{
					Typ:   NumberTT,
					Value: "1",
				},
				ScannedToken{
					Typ:   NumberTT,
					Value: "57.123",
				},
				ScannedToken{
					Typ:   NumberTT,
					Value: "-2",
				},
			},
		},
		{
			Name:  "trailingDecimal",
			Input: `(+ 57. )`,
			Output: []ScannedToken{
				ScannedToken{
					Typ:   OpenParenTT,
					Value: "(",
				},
				ScannedToken{
					Typ:   OpTT,
					Value: "+",
				},
				ScannedToken{
					Typ:   InvalidTT,
					Value: "57.",
				},
			},
		},
		{
			Name:  "BasicOp",
			Input: `(+ 1 234)`,
			Output: []ScannedToken{
				ScannedToken{
					Typ:   OpenParenTT,
					Value: "(",
				},
				ScannedToken{
					Typ:   OpTT,
					Value: "+",
				},
				ScannedToken{
					Typ:   NumberTT,
					Value: "1",
				},
				ScannedToken{
					Typ:   NumberTT,
					Value: "234",
				},
				ScannedToken{
					Typ:   CloseParenTT,
					Value: ")",
				},
			},
		},
		{
			Name:  "basicStr",
			Input: `"abc efg"`,
			Output: []ScannedToken{
				ScannedToken{
					Typ:   StringTT,
					Value: `"abc efg"`,
				},
			},
		},
		{
			Name:  "ident",
			Input: `(let x 5)`,
			Output: []ScannedToken{
				ScannedToken{
					Typ:   OpenParenTT,
					Value: "(",
				},
				ScannedToken{
					Typ:   IdentTT,
					Value: "let",
				},
				ScannedToken{
					Typ:   IdentTT,
					Value: "x",
				},
				ScannedToken{
					Typ:   NumberTT,
					Value: "5",
				},
				ScannedToken{
					Typ:   CloseParenTT,
					Value: ")",
				},
			},
		},
		{
			Name: "comment",
			Input: `
			; 1
			2
			; 3
			4
			`,
			Output: []ScannedToken{
				ScannedToken{
					Typ:   NumberTT,
					Value: "2",
				},
				ScannedToken{
					Typ:   NumberTT,
					Value: "4",
				},
			},
		},
		{
			Name:  "badOperator",
			Input: `--a`,
			Output: []ScannedToken{
				ScannedToken{
					Typ:   InvalidTT,
					Value: "--a",
				},
			},
		},
		{
			Name:  "badNum",
			Input: `123z`,
			Output: []ScannedToken{
				ScannedToken{
					Typ:   InvalidTT,
					Value: "123z",
				},
			},
		},
		{
			Name:  "interruptedString",
			Input: "\"abc\nefg\"",
			Output: []ScannedToken{
				ScannedToken{
					Typ:   InvalidTT,
					Value: "\"abc\n",
				},
			},
		},
		{
			Name:  "smushedString",
			Input: "\"abc\"123",
			Output: []ScannedToken{
				ScannedToken{
					Typ:   InvalidTT,
					Value: "\"abc\"1",
				},
			},
		},
		{
			Name:  "badIdent",
			Input: "abcd++",
			Output: []ScannedToken{
				ScannedToken{
					Typ:   InvalidTT,
					Value: "abcd+",
				},
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			if c.Disabled {
				t.Skip()
			}

			tokens := tokenizeString(fName, c.Input)

			// note (bs): in this case, it might be better to still perform iteration
			// for the shared length, but then error afterwards with whatever
			// missed/expected tokens there are
			if len(tokens) != len(c.Output) {
				t.Fatalf("Token length mismatch \n[expected=%+v]\n[actual=%+v]",
					c.Output, tokens)
			}

			for tokenI, expectedV := range c.Output {
				actualV := tokens[tokenI]

				require.Equalf(t, expectedV.Value, actualV.Value,
					"mismatched values at index %d", tokenI)
				require.Equalf(t, expectedV.Typ.String(), actualV.Typ.String(),
					"mismatched types at index %d", tokenI)
			}
		})
	}

	t.Run("positionTest", func(t *testing.T) {
		actualTokens := tokenizeString(fName, "12\n  34")
		expectedTokens := []ScannedToken{
			ScannedToken{
				Typ:   NumberTT,
				Value: "12",
				Pos:   makePos(1, 1),
			},
			ScannedToken{
				Typ:   NumberTT,
				Value: "34",
				Pos:   makePos(3, 2),
			},
		}
		require.Equal(t, expectedTokens, actualTokens)
	})

	t.Run("invalidChar", func(t *testing.T) {
		actualTokens := tokenizeString(fName, "\x01")
		expectedTokens := []ScannedToken{
			ScannedToken{
				Typ:   InvalidTT,
				Value: "\x01",
				Pos:   makePos(1, 1),
			},
		}
		require.Equal(t, expectedTokens, actualTokens)
	})
}

// tokenizeString converts the provided string to a list of tokens.
func tokenizeString(srcName, str string) []ScannedToken {
	tokens := []ScannedToken{}

	cs := NewRuneScanner(srcName, strings.NewReader(str))
	ts := NewTokenScanner(cs)
	for !ts.Done() {
		ts.Advance()
		nextT := ts.Token()
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

package golisp2

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Tokenization(t *testing.T) {
	fName := "testFile.l"

	testCases := []struct {
		Name     string
		Input    string
		Output   []ScannedToken
		Disabled bool
	}{
		{
			Name:  "OpenClose",
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
			Name:  "Operators",
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
			Name:  "BasicNumbers",
			Input: `1 -2 57.123`,
			Output: []ScannedToken{
				ScannedToken{
					Typ:   NumberTT,
					Value: "1",
				},
				ScannedToken{
					Typ:   NumberTT,
					Value: "-2",
				},
				ScannedToken{
					Typ:   NumberTT,
					Value: "57.123",
				},
			},
		},
		{
			Name:  "TrailingDecimal",
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
			Name:  "BasicStr",
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
					Typ:   CommentTT,
					Value: "; 1",
				},
				ScannedToken{
					Typ:   NumberTT,
					Value: "2",
				},
				ScannedToken{
					Typ:   CommentTT,
					Value: "; 3",
				},
				ScannedToken{
					Typ:   NumberTT,
					Value: "4",
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
				t.Fatalf("Token length mismatch [expected=%+v] [actual=%+v]",
					c.Output, tokens)
			}

			for tokenI, expectedV := range c.Output {
				actualV := tokens[tokenI]

				require.Equalf(t, expectedV.Typ, actualV.Typ,
					"mismatched types at index %d", tokenI)
				require.Equalf(t, expectedV.Value, actualV.Value,
					"mismatched values at index %d", tokenI)
			}
		})
	}

	t.Run("positionTest", func(t *testing.T) {
		makePos := func(c, r int) ScannerPosition {
			return ScannerPosition{
				SourceFile: fName,
				Col:        c,
				Row:        r,
			}
		}

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
}

// tokenizeString converts the provided string to a list of tokens.
func tokenizeString(srcName, str string) []ScannedToken {
	tokens := []ScannedToken{}

	cs := NewRuneScanner(srcName, strings.NewReader(str))
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

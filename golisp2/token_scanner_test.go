package golisp2

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

func Test_TokenizeString(t *testing.T) {
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
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			if c.Disabled {
				t.Skip()
			}

			tokens := TokenizeString(c.Input)

			// note (bs): in this case, it might be better to still perform iteration
			// for the shared length, but then error afterwards with whatever
			// missed/expected tokens there are
			if len(tokens) != len(c.Output) {
				t.Fatalf("Token length mismatch [expected=%+v] [actual=%+v]",
					c.Output, tokens)
			}

			for tokenI, expectedV := range c.Output {
				actualV := tokens[tokenI]
				if expectedV.Typ != actualV.Typ {
					// note (bs): this will be kinda inscrutable as it's not properly
					// "strung"
					t.Fatalf("Mismatched token types at index %d [expected=%s] [actual=%s]",
						tokenI, expectedV.Typ, actualV.Typ)
				}
				if expectedV.Value != actualV.Value {
					t.Fatalf("Mismatched token values at index %d [expected=%s] [actual=%s]",
						tokenI, expectedV.Value, actualV.Value)
				}
			}
		})
	}
}

func Test_abc(t *testing.T) {
	r := strings.NewReader("abcdef")
	var _ = r.Read
	buf := bufio.NewReader(r)
	for i := 0; i < 10; i++ {
		r, n, e := buf.ReadRune()
		var _ = r
		fmt.Println("@@@ values:", n, e)
	}
}

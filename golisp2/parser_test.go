package golisp2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseTokens(t *testing.T) {

	t.Run("basic", func(t *testing.T) {
		assertNumValue(t, evalStrToVal(t, "; hello there\n(+ 1 2)"), 3)
	})

	t.Run("fn", func(t *testing.T) {
		assertNumValue(t, evalStrToVal(t, `((fn (x) (+ x x)) 5)`), 10)
	})

	t.Run("if", func(t *testing.T) {
		assertNumValue(t, evalStrToVal(t, `(if (== 1 2) (+ 5 5) (+ 10 10))`), 20)
	})

	t.Run("str", func(t *testing.T) {
		assertStringValue(t, evalStrToVal(t, `(concat "abc" "efg")`), "abcefg")
	})

	t.Run("bool", func(t *testing.T) {
		assertBoolValue(t, evalStrToVal(t, `(or true false)`), true)
	})

	t.Run("let", func(t *testing.T) {
		assertNumValue(t, evalStrToVal(t, `
		((fn (x)
		  (let y (+ x x))
		  (+ y y))
		 5)`), 20)
	})

	t.Run("operators", func(t *testing.T) {
		t.Run("+", func(t *testing.T) {
			assertNumValue(t, evalStrToVal(t, `(+ 1 1)`), 2)
		})
		t.Run("-", func(t *testing.T) {
			assertNumValue(t, evalStrToVal(t, `(- 1 2)`), -1)
		})
		t.Run("*", func(t *testing.T) {
			assertNumValue(t, evalStrToVal(t, `(* 2 3)`), 6)
		})
		t.Run("/", func(t *testing.T) {
			assertNumValue(t, evalStrToVal(t, `(/ 12 3)`), 4)
		})
		t.Run("==", func(t *testing.T) {
			assertBoolValue(t, evalStrToVal(t, `(== 1 1)`), true)
			assertBoolValue(t, evalStrToVal(t, `(== 1 2)`), false)
			assertBoolValue(t, evalStrToVal(t, `(== 2 1)`), false)
		})
		t.Run("<", func(t *testing.T) {
			assertBoolValue(t, evalStrToVal(t, `(< 1 1)`), false)
			assertBoolValue(t, evalStrToVal(t, `(< 1 2)`), true)
			assertBoolValue(t, evalStrToVal(t, `(< 2 1)`), false)
		})
		t.Run(">", func(t *testing.T) {
			assertBoolValue(t, evalStrToVal(t, `(> 1 1)`), false)
			assertBoolValue(t, evalStrToVal(t, `(> 1 2)`), false)
			assertBoolValue(t, evalStrToVal(t, `(> 2 1)`), true)
		})
		t.Run("<=", func(t *testing.T) {
			assertBoolValue(t, evalStrToVal(t, `(<= 1 1)`), true)
			assertBoolValue(t, evalStrToVal(t, `(<= 1 2)`), true)
			assertBoolValue(t, evalStrToVal(t, `(<= 2 1)`), false)
		})
		t.Run(">=", func(t *testing.T) {
			assertBoolValue(t, evalStrToVal(t, `(>= 1 1)`), true)
			assertBoolValue(t, evalStrToVal(t, `(>= 1 2)`), false)
			assertBoolValue(t, evalStrToVal(t, `(>= 2 1)`), true)
		})
	})

	t.Run("errorsInParse", func(t *testing.T) {

		t.Run("incompleteExpression", func(t *testing.T) {
			parseStrToErr(t, `(+ 1 (+ 2 3`)
		})

		t.Run("invalidToken", func(t *testing.T) {
			parseStrToErr(t, `(+ 1. 2)`)
		})

		t.Run("badOperator", func(t *testing.T) {
			err := parseStrToErr(t, `(++== 1 2)`)
			require.IsType(t, (*ParseError)(nil), err)
			asPE := err.(*ParseError)
			require.Equal(t, "++==", asPE.token.Value)
			require.Equal(t, 2, asPE.token.Pos.Col)
			require.Equal(t, 1, asPE.token.Pos.Row)
		})

		t.Run("invalidFn", func(t *testing.T) {
			parseStrToErr(t, `(fn)`)
			parseStrToErr(t, `(fn (+ 1 2))`)
			parseStrToErr(t, `(fn "abc")`)
			parseStrToErr(t, `(fn (a b 1))`)
		})

		t.Run("invalidLet", func(t *testing.T) {
			parseStrToErr(t, `(let)`)
			parseStrToErr(t, `(let a)`)
			parseStrToErr(t, `(let 1 a)`)
		})

		t.Run("invalidIf", func(t *testing.T) {
			parseStrToErr(t, `(if)`)
		})

		t.Run("misplacedIdent", func(t *testing.T) {
			parseStrToErr(t, `hello`)
		})
	})
}

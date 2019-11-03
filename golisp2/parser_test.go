package golisp2

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseTokens(t *testing.T) {

	evalStringToValue := func(t *testing.T, str string) Value {
		ts := NewTokenScanner(NewRuneScanner(strings.NewReader(str)))
		exprs, exprsErr := ParseTokens(ts)
		require.NoError(t, exprsErr)
		require.Equal(t, len(exprs), 1)
		return mustEval(t, exprs[0], BuiltinContext())
	}

	t.Run("basic", func(t *testing.T) {
		v := evalStringToValue(t, `(+ 1 2)`)
		assertNumValue(t, v, 3)
	})

	t.Run("fn", func(t *testing.T) {
		v := evalStringToValue(t, `((fn (x) (+ x x)) 5)`)
		assertNumValue(t, v, 10)
	})

	t.Run("if", func(t *testing.T) {
		v := evalStringToValue(t, `(if (== 1 2) (+ 5 5) (+ 10 10))`)
		assertNumValue(t, v, 20)
	})

	t.Run("str", func(t *testing.T) {
		v := evalStringToValue(t, `(concat "abc" "efg")`)
		assertStringValue(t, v, "abcefg")
	})

	t.Run("bool", func(t *testing.T) {
		v := evalStringToValue(t, `(or true false)`)
		assertBoolValue(t, v, true)
	})

	t.Run("let", func(t *testing.T) {
		v := evalStringToValue(t, `
		((fn (x)
		  (let y (+ x x))
		  (+ y y))
		 5)`)
		assertNumValue(t, v, 20)
	})
}

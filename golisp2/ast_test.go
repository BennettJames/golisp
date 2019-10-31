package golisp2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_add(t *testing.T) {
	v, e := addFn(
		ExprContext{},
		NewNumberValue(1),
		NewNumberValue(2),
	)
	require.NoError(t, e)
	assertNumValue(t, v, 3)
}

func Test_sub(t *testing.T) {
	v, e := subFn(
		ExprContext{},
		NewNumberValue(2),
		NewNumberValue(1),
	)
	require.NoError(t, e)
	assertNumValue(t, v, 1)
}

func Test_concat(t *testing.T) {
	v, e := concatFn(
		ExprContext{},
		NewStringValue("abc"),
		NewStringValue("efg"),
	)
	require.NoError(t, e)
	assertStringValue(t, v, "abcefg")
}

func Test_cell(t *testing.T) {
	v, e := consFn(
		ExprContext{},
		NewNumberValue(1),
		NewNumberValue(2),
	)
	require.NoError(t, e)
	assertAsCell(t, v)

	left, leftErr := carFn(
		ExprContext{},
		v,
	)
	right, rightErr := cdrFn(
		ExprContext{},
		v,
	)
	require.NoError(t, leftErr)
	require.NoError(t, rightErr)
	assertNumValue(t, left, 1)
	assertNumValue(t, right, 2)
	require.Equal(t, "(1 . 2)", v.PrintStr())
}

func Test_bool(t *testing.T) {
	t.Run("boolPrintStr", func(t *testing.T) {
		require.Equal(t, "true", NewBoolValue(true).PrintStr())
		require.Equal(t, "false", NewBoolValue(false).PrintStr())
	})

	t.Run("and", func(t *testing.T) {
		v1, e1 := andFn(
			ExprContext{},
			NewBoolValue(true),
			NewBoolValue(true),
		)
		require.NoError(t, e1)
		assertBoolValue(t, v1, true)

		v2, e2 := andFn(
			ExprContext{},
			NewBoolValue(true),
			NewBoolValue(false),
		)
		require.NoError(t, e2)
		assertBoolValue(t, v2, false)
	})

	t.Run("or", func(t *testing.T) {
		v1, e1 := orFn(
			ExprContext{},
			NewBoolValue(true),
			NewBoolValue(false),
		)
		require.NoError(t, e1)
		assertBoolValue(t, v1, true)

		v2, e2 := orFn(
			ExprContext{},
			NewBoolValue(false),
			NewBoolValue(false),
			NewBoolValue(false),
		)
		require.NoError(t, e2)
		assertBoolValue(t, v2, false)
	})

	t.Run("not", func(t *testing.T) {
		v1, e1 := notFn(
			ExprContext{},
			NewBoolValue(true),
		)
		require.NoError(t, e1)
		assertBoolValue(t, v1, false)

		v2, e2 := notFn(
			ExprContext{},
			NewBoolValue(false),
		)
		require.NoError(t, e2)
		assertBoolValue(t, v2, true)
	})
}

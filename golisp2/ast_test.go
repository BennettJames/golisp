package golisp2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_mathFn(t *testing.T) {
	t.Run("add", func(t *testing.T) {
		v, e := addFn(
			&ExprContext{},
			NewNumberValue(1),
			NewNumberValue(2),
		)
		require.NoError(t, e)
		assertNumValue(t, v, 3)
	})

	t.Run("sub", func(t *testing.T) {
		v, e := subFn(
			&ExprContext{},
			NewNumberValue(2),
			NewNumberValue(1),
		)
		require.NoError(t, e)
		assertNumValue(t, v, 1)
	})

	t.Run("mult", func(t *testing.T) {
		v, e := multFn(
			&ExprContext{},
			NewNumberValue(3),
			NewNumberValue(4),
		)
		require.NoError(t, e)
		assertNumValue(t, v, 12)
	})

	t.Run("divide", func(t *testing.T) {
		v, e := divFn(
			&ExprContext{},
			NewNumberValue(6),
			NewNumberValue(2),
		)
		require.NoError(t, e)
		assertNumValue(t, v, 3)
	})
}

func Test_concat(t *testing.T) {
	v, e := concatFn(
		&ExprContext{},
		NewStringValue("abc"),
		NewStringValue("efg"),
	)
	require.NoError(t, e)
	assertStringValue(t, v, "abcefg")
}

func Test_cell(t *testing.T) {
	v, e := consFn(
		&ExprContext{},
		NewStringValue("a"),
		NewStringValue("b"),
	)
	require.NoError(t, e)
	assertAsCell(t, v)

	left, leftErr := carFn(
		&ExprContext{},
		v,
	)
	right, rightErr := cdrFn(
		&ExprContext{},
		v,
	)
	require.NoError(t, leftErr)
	require.NoError(t, rightErr)
	assertStringValue(t, left, "a")
	assertStringValue(t, right, "b")
	require.Equal(t, "(\"a\" . \"b\")", v.PrintStr())
}

func Test_bool(t *testing.T) {
	t.Run("boolPrintStr", func(t *testing.T) {
		require.Equal(t, "true", NewBoolValue(true).PrintStr())
		require.Equal(t, "false", NewBoolValue(false).PrintStr())
	})

	t.Run("and", func(t *testing.T) {
		v1, e1 := andFn(
			&ExprContext{},
			NewBoolValue(true),
			NewBoolValue(true),
		)
		require.NoError(t, e1)
		assertBoolValue(t, v1, true)

		v2, e2 := andFn(
			&ExprContext{},
			NewBoolValue(true),
			NewBoolValue(false),
		)
		require.NoError(t, e2)
		assertBoolValue(t, v2, false)
	})

	t.Run("or", func(t *testing.T) {
		v1, e1 := orFn(
			&ExprContext{},
			NewBoolValue(true),
			NewBoolValue(false),
		)
		require.NoError(t, e1)
		assertBoolValue(t, v1, true)

		v2, e2 := orFn(
			&ExprContext{},
			NewBoolValue(false),
			NewBoolValue(false),
			NewBoolValue(false),
		)
		require.NoError(t, e2)
		assertBoolValue(t, v2, false)
	})

	t.Run("not", func(t *testing.T) {
		v1, e1 := notFn(
			&ExprContext{},
			NewBoolValue(true),
		)
		require.NoError(t, e1)
		assertBoolValue(t, v1, false)

		v2, e2 := notFn(
			&ExprContext{},
			NewBoolValue(false),
		)
		require.NoError(t, e2)
		assertBoolValue(t, v2, true)
	})
}

func Test_ident(t *testing.T) {

	ec := &ExprContext{
		parent: &ExprContext{
			vals: map[string]Value{
				"a": NewStringValue("a"),
			},
		},
		vals: map[string]Value{
			"b": NewStringValue("b"),
			"c": NewStringValue("c"),
		},
	}

	v1 := NewIdentValue("a").Eval(ec)
	assertStringValue(t, v1, "a")

	v2 := NewIdentValue("b").Eval(ec)
	assertStringValue(t, v2, "b")

	v3 := NewIdentValue("d").Eval(ec)
	assertNilValue(t, v3)
}

func Test_parenExpr(t *testing.T) {
	ec := &ExprContext{
		parent: &ExprContext{
			vals: map[string]Value{
				"add": NewFuncValue(addFn),
				"sub": NewFuncValue(subFn),
			},
		},
		vals: map[string]Value{
			"a": NewNumberValue(1),
			"b": NewNumberValue(2),
		},
	}
	v := NewCallExpr(
		NewIdentValue("add"),
		NewIdentValue("a"),
		NewIdentValue("b"),
		NewCallExpr(
			NewIdentValue("sub"),
			NewNumberValue(3),
			NewIdentValue("b"),
		),
	).Eval(ec)
	assertNumValue(t, v, 4)
}

func Test_numComparison(t *testing.T) {

	t.Run("eq", func(t *testing.T) {
		v1 := NewCallExpr(
			NewFuncValue(eqNumFn),
			NewNumberValue(1),
			NewNumberValue(1),
		).Eval(nil)
		assertBoolValue(t, v1, true)

		v2 := NewCallExpr(
			NewFuncValue(eqNumFn),
			NewNumberValue(1),
			NewNumberValue(2),
		).Eval(nil)
		assertBoolValue(t, v2, false)

	})

	t.Run("gt", func(t *testing.T) {
		v1 := NewCallExpr(
			NewFuncValue(gtNumFn),
			NewNumberValue(1),
			NewNumberValue(1),
		).Eval(nil)
		assertBoolValue(t, v1, false)

		v2 := NewCallExpr(
			NewFuncValue(gtNumFn),
			NewNumberValue(1),
			NewNumberValue(2),
		).Eval(nil)
		assertBoolValue(t, v2, false)

		v3 := NewCallExpr(
			NewFuncValue(gtNumFn),
			NewNumberValue(2),
			NewNumberValue(1),
		).Eval(nil)
		assertBoolValue(t, v3, true)
	})

	t.Run("lt", func(t *testing.T) {
		v1 := NewCallExpr(
			NewFuncValue(ltNumFn),
			NewNumberValue(1),
			NewNumberValue(1),
		).Eval(nil)
		assertBoolValue(t, v1, false)

		v2 := NewCallExpr(
			NewFuncValue(ltNumFn),
			NewNumberValue(1),
			NewNumberValue(2),
		).Eval(nil)
		assertBoolValue(t, v2, true)

		v3 := NewCallExpr(
			NewFuncValue(ltNumFn),
			NewNumberValue(2),
			NewNumberValue(1),
		).Eval(nil)
		assertBoolValue(t, v3, false)
	})

	t.Run("gte", func(t *testing.T) {
		v1 := NewCallExpr(
			NewFuncValue(gteNumFn),
			NewNumberValue(1),
			NewNumberValue(1),
		).Eval(nil)
		assertBoolValue(t, v1, true)

		v2 := NewCallExpr(
			NewFuncValue(gteNumFn),
			NewNumberValue(1),
			NewNumberValue(2),
		).Eval(nil)
		assertBoolValue(t, v2, false)

		v3 := NewCallExpr(
			NewFuncValue(gteNumFn),
			NewNumberValue(2),
			NewNumberValue(1),
		).Eval(nil)
		assertBoolValue(t, v3, true)
	})

	t.Run("lte", func(t *testing.T) {
		v1 := NewCallExpr(
			NewFuncValue(lteNumFn),
			NewNumberValue(1),
			NewNumberValue(1),
		).Eval(nil)
		assertBoolValue(t, v1, true)

		v2 := NewCallExpr(
			NewFuncValue(lteNumFn),
			NewNumberValue(1),
			NewNumberValue(2),
		).Eval(nil)
		assertBoolValue(t, v2, true)

		v3 := NewCallExpr(
			NewFuncValue(lteNumFn),
			NewNumberValue(2),
			NewNumberValue(1),
		).Eval(nil)
		assertBoolValue(t, v3, false)
	})
}

func Test_ifExpr(t *testing.T) {
	v1 := NewIfExpr(
		NewBoolValue(true),
		NewNumberValue(1),
		NewNumberValue(2),
	).Eval(nil)
	assertNumValue(t, v1, 1)
	v2 := NewIfExpr(
		NewBoolValue(false),
		NewNumberValue(1),
		NewNumberValue(2),
	).Eval(nil)
	assertNumValue(t, v2, 2)
}

func Test_fnExpr(t *testing.T) {

	doubleAdd := NewFnExpr(
		[]Arg{
			Arg{
				Ident: "a",
			},
			Arg{
				Ident: "b",
			},
		},
		[]Expr{
			NewCallExpr(
				NewFuncValue(addFn),
				NewIdentValue("a"),
				NewIdentValue("b"),
				NewIdentValue("b"),
			),
		},
	).Eval(nil)
	asFn := assertAsFunc(t, doubleAdd)

	v, e := asFn.Exec(nil, NewNumberValue(1), NewNumberValue(2))
	require.NoError(t, e)
	assertNumValue(t, v, 5)
}

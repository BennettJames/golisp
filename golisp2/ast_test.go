package golisp2

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_mathFn(t *testing.T) {
	t.Run("add", func(t *testing.T) {
		v, e := addFn(
			BuiltinContext(),
			NewNumberValue(1),
			NewNumberValue(2),
		)
		require.NoError(t, e)
		assertNumValue(t, v, 3)
	})

	t.Run("sub", func(t *testing.T) {
		v, e := subFn(
			BuiltinContext(),
			NewNumberValue(2),
			NewNumberValue(1),
		)
		require.NoError(t, e)
		assertNumValue(t, v, 1)
	})

	t.Run("mult", func(t *testing.T) {
		v, e := multFn(
			BuiltinContext(),
			NewNumberValue(3),
			NewNumberValue(4),
		)
		require.NoError(t, e)
		assertNumValue(t, v, 12)
	})

	t.Run("divide", func(t *testing.T) {
		v, e := divFn(
			BuiltinContext(),
			NewNumberValue(6),
			NewNumberValue(2),
		)
		require.NoError(t, e)
		assertNumValue(t, v, 3)
	})
}

func Test_concat(t *testing.T) {
	v, e := concatFn(
		BuiltinContext(),
		NewStringValue("abc"),
		NewStringValue("efg"),
	)
	require.NoError(t, e)
	assertStringValue(t, v, "abcefg")
}

func Test_cell(t *testing.T) {
	v, e := consFn(
		BuiltinContext(),
		NewStringValue("a"),
		NewStringValue("b"),
	)
	require.NoError(t, e)
	assertAsCell(t, v)

	left, leftErr := carFn(
		BuiltinContext(),
		v,
	)
	right, rightErr := cdrFn(
		BuiltinContext(),
		v,
	)
	require.NoError(t, leftErr)
	require.NoError(t, rightErr)
	assertStringValue(t, left, "a")
	assertStringValue(t, right, "b")
	require.Equal(t, "(\"a\" . \"b\")", v.InspectStr())
}

func Test_bool(t *testing.T) {
	t.Run("boolInspectStr", func(t *testing.T) {
		require.Equal(t, "true", NewBoolValue(true).InspectStr())
		require.Equal(t, "false", NewBoolValue(false).InspectStr())
	})

	t.Run("and", func(t *testing.T) {
		v1, e1 := andFn(
			BuiltinContext(),
			NewBoolValue(true),
			NewBoolValue(true),
		)
		require.NoError(t, e1)
		assertBoolValue(t, v1, true)

		v2, e2 := andFn(
			BuiltinContext(),
			NewBoolValue(true),
			NewBoolValue(false),
		)
		require.NoError(t, e2)
		assertBoolValue(t, v2, false)
	})

	t.Run("or", func(t *testing.T) {
		v1, e1 := orFn(
			BuiltinContext(),
			NewBoolValue(true),
			NewBoolValue(false),
		)
		require.NoError(t, e1)
		assertBoolValue(t, v1, true)

		v2, e2 := orFn(
			BuiltinContext(),
			NewBoolValue(false),
			NewBoolValue(false),
			NewBoolValue(false),
		)
		require.NoError(t, e2)
		assertBoolValue(t, v2, false)
	})

	t.Run("not", func(t *testing.T) {
		v1, e1 := notFn(
			BuiltinContext(),
			NewBoolValue(true),
		)
		require.NoError(t, e1)
		assertBoolValue(t, v1, false)

		v2, e2 := notFn(
			BuiltinContext(),
			NewBoolValue(false),
		)
		require.NoError(t, e2)
		assertBoolValue(t, v2, true)
	})
}

func Test_ident(t *testing.T) {

	ec := &EvalContext{
		parent: &EvalContext{
			vals: map[string]Value{
				"a": NewStringValue("a"),
			},
		},
		vals: map[string]Value{
			"b": NewStringValue("b"),
			"c": NewStringValue("c"),
		},
	}

	v1 := mustEval(t, NewIdentValue("a"), ec)
	assertStringValue(t, v1, "a")

	v2 := mustEval(t, NewIdentValue("b"), ec)
	assertStringValue(t, v2, "b")

	v3 := mustEval(t, NewIdentValue("d"), ec)
	assertNilValue(t, v3)
}

func Test_parenExpr(t *testing.T) {
	ec := &EvalContext{
		parent: &EvalContext{
			vals: map[string]Value{
				"add": NewFuncValue("", addFn),
				"sub": NewFuncValue("", subFn),
			},
		},
		vals: map[string]Value{
			"a": NewNumberValue(1),
			"b": NewNumberValue(2),
		},
	}
	v := mustEval(t,
		NewCallExpr(
			NewIdentValue("add"),
			NewIdentValue("a"),
			NewIdentValue("b"),
			NewCallExpr(
				NewIdentValue("sub"),
				NewNumberValue(3),
				NewIdentValue("b"),
			),
		),
		ec,
	)
	assertNumValue(t, v, 4)
}

func Test_numComparison(t *testing.T) {

	t.Run("eq", func(t *testing.T) {
		v1 := mustEval(t,
			NewCallExpr(
				NewFuncValue("", eqNumFn),
				NewNumberValue(1),
				NewNumberValue(1),
			),
			nil)
		assertBoolValue(t, v1, true)

		v2 := mustEval(t,
			NewCallExpr(
				NewFuncValue("", eqNumFn),
				NewNumberValue(1),
				NewNumberValue(2),
			), nil)
		assertBoolValue(t, v2, false)

	})

	t.Run("gt", func(t *testing.T) {
		v1 := mustEval(t, NewCallExpr(
			NewFuncValue("", gtNumFn),
			NewNumberValue(1),
			NewNumberValue(1),
		), nil)
		assertBoolValue(t, v1, false)

		v2 := mustEval(t, NewCallExpr(
			NewFuncValue("", gtNumFn),
			NewNumberValue(1),
			NewNumberValue(2),
		), nil)
		assertBoolValue(t, v2, false)

		v3 := mustEval(t, NewCallExpr(
			NewFuncValue("", gtNumFn),
			NewNumberValue(2),
			NewNumberValue(1),
		), nil)
		assertBoolValue(t, v3, true)
	})

	t.Run("lt", func(t *testing.T) {
		v1 := mustEval(t, NewCallExpr(
			NewFuncValue("", ltNumFn),
			NewNumberValue(1),
			NewNumberValue(1),
		), nil)
		assertBoolValue(t, v1, false)

		v2 := mustEval(t, NewCallExpr(
			NewFuncValue("", ltNumFn),
			NewNumberValue(1),
			NewNumberValue(2),
		), nil)
		assertBoolValue(t, v2, true)

		v3 := mustEval(t, NewCallExpr(
			NewFuncValue("", ltNumFn),
			NewNumberValue(2),
			NewNumberValue(1),
		), nil)
		assertBoolValue(t, v3, false)
	})

	t.Run("gte", func(t *testing.T) {
		v1 := mustEval(t, NewCallExpr(
			NewFuncValue("", gteNumFn),
			NewNumberValue(1),
			NewNumberValue(1),
		), nil)
		assertBoolValue(t, v1, true)

		v2 := mustEval(t, NewCallExpr(
			NewFuncValue("", gteNumFn),
			NewNumberValue(1),
			NewNumberValue(2),
		), nil)
		assertBoolValue(t, v2, false)

		v3 := mustEval(t, NewCallExpr(
			NewFuncValue("", gteNumFn),
			NewNumberValue(2),
			NewNumberValue(1),
		), nil)
		assertBoolValue(t, v3, true)
	})

	t.Run("lte", func(t *testing.T) {
		v1 := mustEval(t, NewCallExpr(
			NewFuncValue("", lteNumFn),
			NewNumberValue(1),
			NewNumberValue(1),
		), nil)
		assertBoolValue(t, v1, true)

		v2 := mustEval(t, NewCallExpr(
			NewFuncValue("", lteNumFn),
			NewNumberValue(1),
			NewNumberValue(2),
		), nil)
		assertBoolValue(t, v2, true)

		v3 := mustEval(t, NewCallExpr(
			NewFuncValue("", lteNumFn),
			NewNumberValue(2),
			NewNumberValue(1),
		), nil)
		assertBoolValue(t, v3, false)
	})
}

func Test_ifExpr(t *testing.T) {
	v1 := mustEval(t, NewIfExpr(
		NewBoolValue(true),
		NewNumberValue(1),
		NewNumberValue(2),
	), nil)
	assertNumValue(t, v1, 1)
	v2 := mustEval(t, NewIfExpr(
		NewBoolValue(false),
		NewNumberValue(1),
		NewNumberValue(2),
	), nil)
	assertNumValue(t, v2, 2)
}

func Test_fnExpr(t *testing.T) {

	doubleAdd := mustEval(t, NewFnExpr(
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
				NewFuncValue("", addFn),
				NewIdentValue("a"),
				NewIdentValue("b"),
				NewIdentValue("b"),
			),
		},
	), nil)
	asFn := assertAsFunc(t, doubleAdd)

	v, e := asFn.Exec(nil, NewNumberValue(1), NewNumberValue(2))
	require.NoError(t, e)
	assertNumValue(t, v, 5)
}

func Test_CodeStr(t *testing.T) {

	// printAndReparse is a helper that converts the expression to string, parses
	// it, and returns the re-parsed expression. Anything unexpected will result
	// in a test failure.
	printAndReparse := func(t *testing.T, e Expr) Expr {
		r := strings.NewReader(e.CodeStr())
		ts := NewTokenScanner(NewRuneScanner("testfile", r))
		exprs, exprsErr := ParseTokens(ts)
		require.NoError(t, exprsErr)
		require.Equal(t, 1, len(exprs))
		return exprs[0]
	}

	t.Run("cell", func(t *testing.T) {
		baseAST := NewCellValue(
			NewNumberValue(1),
			NewCellValue(
				NewNumberValue(2),
				nil,
			),
		)
		reparsedExpr := printAndReparse(t, baseAST)
		require.Equal(t, baseAST, mustEval(t, reparsedExpr, nil))
	})

	t.Run("if", func(t *testing.T) {
		baseAST := &IfExpr{
			Cond: NewBoolValue(false),
			Case1: NewCallExpr(
				NewIdentValue("car"),
				NewCellValue(
					NewNumberValue(1),
					NewNumberValue(2),
				),
			),
			Case2: NewCallExpr(
				NewIdentValue("cdr"),
				NewCellValue(
					NewNumberValue(1),
					NewNumberValue(2),
				),
			),
		}
		reparsedExpr := printAndReparse(t, baseAST)
		assertNumValue(t, mustEval(t, reparsedExpr, nil), 2)
	})

	t.Run("let", func(t *testing.T) {
		baseAST := &LetExpr{
			Ident: NewIdentValue("value"),
			Value: NewNumberValue(2),
		}
		reparsedExpr := printAndReparse(t, baseAST)
		ec := BuiltinContext()
		reparsedExpr.Eval(ec)
		ctxVal, hasCtxVal := ec.Resolve("value")
		require.True(t, hasCtxVal)
		assertNumValue(t, ctxVal, 2)
	})

	t.Run("fn", func(t *testing.T) {
		baseAST := NewCallExpr(
			NewFnExpr(
				[]Arg{
					{Ident: "a"},
				},
				[]Expr{
					NewCallExpr(
						NewIdentValue("add"),
						NewIdentValue("a"),
						NewNumberValue(1),
					),
				},
			),
			NewNumberValue(5),
		)
		reparsedExpr := printAndReparse(t, baseAST)
		v := mustEval(t, reparsedExpr, BuiltinContext().SubContext(map[string]Value{
			"add": NewFuncValue("add", addFn),
		}))
		assertNumValue(t, v, 6)
	})
}

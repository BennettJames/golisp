package golisp2

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_cell(t *testing.T) {
	v, e := NewCallExpr(
		NewIdentLiteral("cons"),
		NewStringLiteral("a"),
		NewStringLiteral("b"),
	).Eval(BuiltinContext())

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

func Test_ident(t *testing.T) {

	ec := &EvalContext{
		parent: &EvalContext{
			vals: map[string]Value{
				"a": &StringValue{
					Val: "a",
				},
			},
		},
		vals: map[string]Value{
			"b": &StringValue{
				Val: "b",
			},
			"c": &StringValue{
				Val: "c",
			},
		},
	}

	v1 := mustEval(t, NewIdentLiteral("a"), ec)
	assertStringValue(t, v1, "a")

	v2 := mustEval(t, NewIdentLiteral("b"), ec)
	assertStringValue(t, v2, "b")

	v3 := mustEval(t, NewIdentLiteral("d"), ec)
	assertNilValue(t, v3)
}

func Test_parenExpr(t *testing.T) {
	ec := &EvalContext{
		parent: &EvalContext{
			vals: map[string]Value{
				"add": &FuncValue{
					Fn: addFn,
				},
				"sub": &FuncValue{
					Fn: subFn,
				},
			},
		},
		vals: map[string]Value{
			"a": &NumberValue{
				Val: 1,
			},
			"b": &NumberValue{
				Val: 2,
			},
		},
	}
	v := mustEval(t,
		NewCallExpr(
			NewIdentLiteral("add"),
			NewIdentLiteral("a"),
			NewIdentLiteral("b"),
			NewCallExpr(
				NewIdentLiteral("sub"),
				NewNumberLiteral(3),
				NewIdentLiteral("b"),
			),
		),
		ec,
	)
	assertNumValue(t, v, 4)
}

func Test_ifExpr(t *testing.T) {

	t.Run("basic", func(t *testing.T) {
		v1 := mustEval(t, NewIfExpr(
			NewBoolLiteral(true),
			NewNumberLiteral(1),
			NewNumberLiteral(2),
		), nil)
		assertNumValue(t, v1, 1)
		v2 := mustEval(t, NewIfExpr(
			NewBoolLiteral(false),
			NewNumberLiteral(1),
			NewNumberLiteral(2),
		), nil)
		assertNumValue(t, v2, 2)
	})

	t.Run("errors", func(t *testing.T) {
		v, err := NewIfExpr(
			NewNumberLiteral(0),
			NewNumberLiteral(1),
			NewNumberLiteral(2),
		).Eval(BuiltinContext())
		require.Error(t, err)
		require.IsType(t, (*TypeError)(nil), err)
		require.Nil(t, v)
	})
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
				NewFuncLiteral("", addFn),
				NewIdentLiteral("a"),
				NewIdentLiteral("b"),
				NewIdentLiteral("b"),
			),
		},
	), nil)
	asFn := assertAsFunc(t, doubleAdd)

	v, e := asFn.Fn(nil, &NumberValue{Val: 1}, &NumberValue{Val: 2})
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
		baseAST := NewCallExpr(
			NewIdentLiteral("car"),
			NewCallExpr(
				NewIdentLiteral("cons"),
				NewNumberLiteral(1),
				NewCallExpr(
					NewIdentLiteral("cons"),
					NewNumberLiteral(2),
					NewNilLiteral(),
				),
			),
		)
		reparsedExpr := printAndReparse(t, baseAST)
		assertNumValue(t, mustEval(t, reparsedExpr, nil), 1)
	})

	t.Run("if", func(t *testing.T) {
		baseAST := &IfExpr{
			Cond: NewBoolLiteral(false),
			Case1: NewCallExpr(
				NewIdentLiteral("car"),

				NewCallExpr(
					NewIdentLiteral("cons"),
					NewNumberLiteral(1),
					NewNumberLiteral(2),
				),
			),
			Case2: NewCallExpr(
				NewIdentLiteral("cdr"),
				NewCallExpr(
					NewIdentLiteral("cons"),
					NewNumberLiteral(1),
					NewNumberLiteral(2),
				),
			),
		}
		reparsedExpr := printAndReparse(t, baseAST)
		assertNumValue(t, mustEval(t, reparsedExpr, nil), 2)
	})

	t.Run("let", func(t *testing.T) {
		baseAST := &LetExpr{
			Ident: NewIdentLiteral("value"),
			Value: NewNumberLiteral(2),
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
					{Ident: "b"},
				},
				[]Expr{
					NewCallExpr(
						NewIdentLiteral("add"),
						NewIdentLiteral("a"),
						NewIdentLiteral("b"),
					),
				},
			),
			NewNumberLiteral(5),
			NewNumberLiteral(1),
		)
		reparsedExpr := printAndReparse(t, baseAST)
		v := mustEval(t, reparsedExpr, BuiltinContext().SubContext(map[string]Value{
			"add": &FuncValue{Fn: addFn},
		}))
		assertNumValue(t, v, 6)
	})
}

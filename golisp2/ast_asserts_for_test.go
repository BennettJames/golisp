package golisp2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func mustEval(t *testing.T, e Expr, ec *ExprContext) Value {
	t.Helper()
	if ec == nil {
		ec = BuiltinContext()
	}
	v, err := e.Eval(ec)
	require.NoError(t, err)
	return v
}

func assertAsNum(t *testing.T, v Value) *NumberValue {
	t.Helper()
	require.NotNil(t, v)
	asNum, isNum := v.(*NumberValue)
	require.True(t, isNum)
	return asNum
}

func assertNumValue(t *testing.T, v Value, expected float64) {
	t.Helper()
	asNum := assertAsNum(t, v)
	require.Equal(t, expected, asNum.Val)
}

func assertNilValue(t *testing.T, v Value) {
	t.Helper()
	require.NotNil(t, v)
	_, isNil := v.(*NilValue)
	require.True(t, isNil)
}

func assertAsString(t *testing.T, v Value) *StringValue {
	t.Helper()
	require.NotNil(t, v)
	asStr, isStr := v.(*StringValue)
	require.True(t, isStr)
	return asStr
}

func assertStringValue(t *testing.T, v Value, expected string) {
	t.Helper()
	asStr := assertAsString(t, v)
	require.Equal(t, expected, asStr.Val)
}

func assertAsBool(t *testing.T, v Value) *BoolValue {
	t.Helper()
	require.NotNil(t, v)
	asBool, isBool := v.(*BoolValue)
	require.True(t, isBool)
	return asBool
}

func assertBoolValue(t *testing.T, v Value, expected bool) {
	t.Helper()
	asBool := assertAsBool(t, v)
	require.Equal(t, expected, asBool.Val)
}

func assertAsFunc(t *testing.T, v Value) *FuncValue {
	t.Helper()
	require.NotNil(t, v)
	asFunc, isFunc := v.(*FuncValue)
	require.True(t, isFunc)
	return asFunc
}

func assertAsCell(t *testing.T, v Value) *CellValue {
	t.Helper()
	require.NotNil(t, v)
	asCell, isCell := v.(*CellValue)
	require.True(t, isCell)
	return asCell
}

func assertCellValue(t *testing.T, v Value, expectedL, expectedR Value) {
	t.Helper()
	asCell := assertAsCell(t, v)
	// note (bs): not 100% convinced this will work well. Let's play around and
	// see if it is sane enough to be useful.
	require.EqualValues(t, expectedL, asCell.Left, "left values should be equal")
	require.EqualValues(t, expectedR, asCell.Right, "right values should be equal")
}

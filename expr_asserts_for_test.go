package golisp2

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// mustEval executes the given expression in the context, and asserts there is
// no error. Returns the subsequent value.
func mustEval(t *testing.T, e Expr, ec *EvalContext) Value {
	t.Helper()
	if ec == nil {
		ec = BuiltinContext()
	}
	v, err := e.Eval(ec)
	require.NoError(t, err)
	return v
}

// evalStrToVal will parse the string, assert that exactly one expression is
// returned, evaluate it and return the result.
func evalStrToVal(t *testing.T, str string) Value {
	t.Helper()
	ts := NewTokenScanner(NewRuneScanner("testfile", strings.NewReader(str)))
	exprs, exprsErr := ParseTokens(ts)
	require.NoError(t, exprsErr)
	require.Equal(t, len(exprs), 1)
	return mustEval(t, exprs[0], BuiltinContext())
}

// evalStrToVal will parse the string, assert that exactly one expression is
// returned, evaluate it, assert it's an error, and return it.
func evalStrToErr(t *testing.T, str string) error {
	t.Helper()
	ts := NewTokenScanner(NewRuneScanner("testfile", strings.NewReader(str)))
	exprs, exprsErr := ParseTokens(ts)
	require.NoError(t, exprsErr)
	require.Equal(t, len(exprs), 1)
	_, err := exprs[0].Eval(BuiltinContext())
	require.Error(t, err)
	return err
}

// parseStrToErr will parse the string, and assert that it results in an error.
func parseStrToErr(t *testing.T, str string) error {
	t.Helper()
	ts := NewTokenScanner(NewRuneScanner("testfile", strings.NewReader(str)))
	_, exprsErr := ParseTokens(ts)
	require.Error(t, exprsErr)
	return exprsErr
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

func assertAsList(t *testing.T, v Value) *ListValue {
	t.Helper()
	require.NotNil(t, v)
	asList, isList := v.(*ListValue)
	require.True(t, isList)
	return asList
}

func assertListValue(t *testing.T, actual Value, expected []Value) {
	t.Helper()
	asList := assertAsList(t, actual)
	// note (bs): not sure if require is smart enough for this; may need to
	// eventually add a more sensitive notion of collection equality.
	require.EqualValues(t, expected, asList.Vals, "list values should be equal")
}

func assertAsMap(t *testing.T, v Value) *MapValue {
	t.Helper()
	require.NotNil(t, v)
	asMap, isMap := v.(*MapValue)
	require.True(t, isMap)
	return asMap
}

func assertMapValue(t *testing.T, actual Value, expected map[string]Value) {
	t.Helper()
	asMap := assertAsMap(t, actual)
	// note (bs): not sure if require is smart enough for this; may need to
	// eventually add a more sensitive notion of collection equality.
	require.EqualValues(t, expected, asMap.Vals, "map values should be equal")
}

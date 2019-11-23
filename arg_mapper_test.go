package golisp2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ArgMapper(t *testing.T) {

	t.Run("basic", func(t *testing.T) {
		args := []Value{
			&NumberValue{Val: 1},
			&StringValue{Val: "abc"},
			&BoolValue{Val: true},
			&FuncValue{Fn: consFn},
			&CellValue{Left: &NilValue{}, Right: &NilValue{}},
			&ListValue{Vals: []Value{&NilValue{}}},
			&MapValue{Vals: map[string]Value{"a": &NilValue{}}},
		}

		var nv *NumberValue
		var sv *StringValue
		var bv *BoolValue
		var fv *FuncValue
		var cv *CellValue
		var lv *ListValue
		var mv *MapValue

		mapErr := ArgMapperValues(args...).
			ReadNumber(&nv).
			ReadString(&sv).
			ReadBool(&bv).
			ReadFunc(&fv).
			ReadCell(&cv).
			ReadList(&lv).
			ReadMap(&mv).
			Err()
		require.NoError(t, mapErr)

		require.NotNil(t, nv)
		require.Equal(t, 1.0, nv.Val)
		require.NotNil(t, sv)
		require.Equal(t, "abc", sv.Val)
		require.NotNil(t, bv)
		require.Equal(t, true, bv.Val)
		require.NotNil(t, fv)
		require.NotNil(t, cv)
		require.NotNil(t, lv)
		require.Equal(t, 1, len(lv.Vals))
		require.NotNil(t, mv)
		require.Equal(t, 1, len(mv.Vals))
	})

	t.Run("tooManyReads", func(t *testing.T) {
		args := []Value{
			&NumberValue{Val: 1},
		}

		var nv *NumberValue
		var sv *StringValue
		var bv *BoolValue

		mapErr := ArgMapperValues(args...).
			ReadNumber(&nv).
			ReadString(&sv).
			ReadBool(&bv).
			Err()
		require.Error(t, mapErr)

		require.NotNil(t, nv)
		require.Equal(t, 1.0, nv.Val)
		require.Nil(t, sv)
		require.Nil(t, bv)
	})

	t.Run("exprMapper", func(t *testing.T) {
		ec := BuiltinContext()
		args := []Expr{
			&NumberLiteral{Num: 1},
			NewCallExpr(
				&FuncLiteral{Fn: addFn},
				&NumberLiteral{Num: 1},
				&NilLiteral{},
			),
		}

		var nv1 *NumberValue
		var nv2 *NumberValue

		mapErr := ArgMapperExprs(ec, args).
			ReadNumber(&nv1).
			ReadNumber(&nv2).
			Err()
		require.Error(t, mapErr, "expr mapper should carry evaluation errors")

		require.NotNil(t, nv1)
		require.Equal(t, 1.0, nv1.Val)
		require.Nil(t, nv2)
	})
}

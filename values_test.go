package golisp2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_boolValue(t *testing.T) {
	t.Run("InspectStr", func(t *testing.T) {
		b1 := &BoolValue{
			Val: true,
		}
		b2 := &BoolValue{
			Val: false,
		}
		require.Equal(t, "true", b1.InspectStr())
		require.Equal(t, "false", b2.InspectStr())
	})
}

func Test_listValue(t *testing.T) {

	t.Run("create", func(t *testing.T) {
		assertListValue(
			t,
			evalStrToVal(t, `(list 1 2 3)`),
			[]Value{
				&NumberValue{1},
				&NumberValue{2},
				&NumberValue{3},
			},
		)
	})

	t.Run("inspect", func(t *testing.T) {
		require.Equal(
			t,
			`["a" "b" "c"]`,
			(&ListValue{
				Vals: []Value{
					&StringValue{Val: "a"},
					&StringValue{Val: "b"},
					&StringValue{Val: "c"},
				},
			}).InspectStr(),
		)
	})

	t.Run("filter", func(t *testing.T) {
		t.Run("basic", func(t *testing.T) {
			assertListValue(
				t,
				evalStrToVal(t, `(filter (list 1 2 3) (fn (v) (== v 2)))`),
				[]Value{
					&NumberValue{2},
				},
			)
		})

		t.Run("allowsNils", func(t *testing.T) {
			assertListValue(
				t,
				evalStrToVal(t, `(filter (list 1 2 3) (fn (v) nil))`),
				[]Value{},
			)
		})

		t.Run("badArgCount", func(t *testing.T) {
			evalStrToErr(t, `(filter (list 1 2 3))`)
		})

		t.Run("badValue", func(t *testing.T) {
			evalStrToErr(t, `(filter (list 1 nil 3) (fn (v) (== v 2)))`)
		})

		t.Run("badReturnValue", func(t *testing.T) {
			evalStrToErr(t, `(filter (list 1 2 3) (fn (v) (+ v 1)))`)
		})

		t.Run("badList", func(t *testing.T) {
			evalStrToErr(t, `(filter "" (fn (v) (== v 2)))`)
		})

		t.Run("badFn", func(t *testing.T) {
			evalStrToErr(t, `(filter (list 1 2 3) "")`)
		})
	})

	t.Run("map", func(t *testing.T) {
		t.Run("basic", func(t *testing.T) {
			assertListValue(
				t,
				evalStrToVal(t, `(map (list 1 2 3) (fn (v) (+ v 1)))`),
				[]Value{
					&NumberValue{2},
					&NumberValue{3},
					&NumberValue{4},
				},
			)
		})

		t.Run("badArgCount", func(t *testing.T) {
			evalStrToErr(t, `(map (list 1 2 3))`)
		})

		t.Run("badValue", func(t *testing.T) {
			evalStrToErr(t, `(map (list 1 nil 3) (fn (v) (+ v 1)))`)
		})

		t.Run("badList", func(t *testing.T) {
			evalStrToErr(t, `(map "hello there" (fn (v) (+ v 1)))`)
		})

		t.Run("badFn", func(t *testing.T) {
			evalStrToErr(t, `(map (list 1 2 3) "hello there")`)
		})
	})

	t.Run("reduce", func(t *testing.T) {
		t.Run("basic", func(t *testing.T) {
			assertNumValue(
				t,
				evalStrToVal(t, `(reduce 1 (list 1 2 3) (fn (t v) (+ t v)))`),
				7.0,
			)
		})

		t.Run("badArgCount", func(t *testing.T) {
			evalStrToErr(t, `(reduce 1 (list 1 2 3))`)
		})

		t.Run("badValue", func(t *testing.T) {
			evalStrToErr(t, `(reduce 1 (list 1 nil 3) (fn (t v) (+ t v)))`)
		})

		t.Run("badList", func(t *testing.T) {
			evalStrToErr(t, `(reduce 1 "hello there" (fn (t v) (+ t v)))`)
		})

		t.Run("badFn", func(t *testing.T) {
			evalStrToErr(t, `(reduce 1 (list 1 2 3) "hello there")`)
		})
	})
}

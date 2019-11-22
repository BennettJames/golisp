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
				evalStrToVal(t, `(listFilter (list 1 2 3) (fn (v) (== v 2)))`),
				[]Value{
					&NumberValue{2},
				},
			)
		})

		t.Run("allowsNils", func(t *testing.T) {
			assertListValue(
				t,
				evalStrToVal(t, `(listFilter (list 1 2 3) (fn (v) nil))`),
				[]Value{},
			)
		})

		t.Run("badArgCount", func(t *testing.T) {
			evalStrToErr(t, `(listFilter (list 1 2 3))`)
		})

		t.Run("badValue", func(t *testing.T) {
			evalStrToErr(t, `(listFilter (list 1 nil 3) (fn (v) (== v 2)))`)
		})

		t.Run("badReturnValue", func(t *testing.T) {
			evalStrToErr(t, `(listFilter (list 1 2 3) (fn (v) (+ v 1)))`)
		})

		t.Run("badList", func(t *testing.T) {
			evalStrToErr(t, `(listFilter "" (fn (v) (== v 2)))`)
		})

		t.Run("badFn", func(t *testing.T) {
			evalStrToErr(t, `(listFilter (list 1 2 3) "")`)
		})
	})

	t.Run("map", func(t *testing.T) {
		t.Run("basic", func(t *testing.T) {
			assertListValue(
				t,
				evalStrToVal(t, `(listMap (list 1 2 3) (fn (v) (+ v 1)))`),
				[]Value{
					&NumberValue{2},
					&NumberValue{3},
					&NumberValue{4},
				},
			)
		})

		t.Run("badArgCount", func(t *testing.T) {
			evalStrToErr(t, `(listMap (list 1 2 3))`)
		})

		t.Run("badValue", func(t *testing.T) {
			evalStrToErr(t, `(listMap (list 1 nil 3) (fn (v) (+ v 1)))`)
		})

		t.Run("badList", func(t *testing.T) {
			evalStrToErr(t, `(listMap "hello there" (fn (v) (+ v 1)))`)
		})

		t.Run("badFn", func(t *testing.T) {
			evalStrToErr(t, `(listMap (list 1 2 3) "hello there")`)
		})
	})

	t.Run("reduce", func(t *testing.T) {
		t.Run("basic", func(t *testing.T) {
			assertNumValue(
				t,
				evalStrToVal(t, `(listReduce 1 (list 1 2 3) (fn (t v) (+ t v)))`),
				7.0,
			)
		})

		t.Run("badArgCount", func(t *testing.T) {
			evalStrToErr(t, `(listReduce 1 (list 1 2 3))`)
		})

		t.Run("badValue", func(t *testing.T) {
			evalStrToErr(t, `(listReduce 1 (list 1 nil 3) (fn (t v) (+ t v)))`)
		})

		t.Run("badList", func(t *testing.T) {
			evalStrToErr(t, `(listReduce 1 "hello there" (fn (t v) (+ t v)))`)
		})

		t.Run("badFn", func(t *testing.T) {
			evalStrToErr(t, `(listReduce 1 (list 1 2 3) "hello there")`)
		})
	})
}

func Test_mapValue(t *testing.T) {

	t.Run("create", func(t *testing.T) {
		assertMapValue(
			t,
			evalStrToVal(t, `(map "a" 1 "b" 2)`),
			map[string]Value{
				"a": &NumberValue{Val: 1},
				"b": &NumberValue{Val: 2},
			},
		)
	})

	t.Run("badCreate", func(t *testing.T) {
		evalStrToErr(t, `(map "a" 1 "b")`)
	})

	t.Run("badKey", func(t *testing.T) {
		evalStrToErr(t, `(map "a" 1 "b" 2 (list 1 2 3) 3)`)
	})

	t.Run("inspectStr", func(t *testing.T) {
		t.Run("stringKey", func(t *testing.T) {
			require.Equal(
				t,
				`{ a:true }`,
				(&MapValue{
					Vals: map[string]Value{
						"a": &BoolValue{Val: true},
					},
				}).InspectStr(),
			)
		})
	})

	t.Run("mapKeys", func(t *testing.T) {
		t.Run("basic", func(t *testing.T) {
			require.ElementsMatch(
				t,
				[]Value{
					&StringValue{Val: "a"},
					&StringValue{Val: "b"},
				},
				assertAsList(t, evalStrToVal(t, `(mapKeys (map "a" 1 "b" 2))`)).Vals,
			)
		})

		t.Run("badArg", func(t *testing.T) {
			evalStrToErr(t, `(mapKeys (list 1 2 3))`)
		})

		t.Run("badArgCount", func(t *testing.T) {
			evalStrToErr(t, `(mapKeys (map "a" 1 "b" 2) (map "a" 1 "b" 2))`)
		})
	})

	t.Run("mapValues", func(t *testing.T) {
		t.Run("basic", func(t *testing.T) {
			require.ElementsMatch(
				t,
				[]Value{
					&NumberValue{Val: 1},
					&NumberValue{Val: 2},
				},
				assertAsList(t, evalStrToVal(t, `(mapValues (map "a" 1 "b" 2))`)).Vals,
			)
		})

		t.Run("badArg", func(t *testing.T) {
			evalStrToErr(t, `(mapValues (list 1 2 3))`)
		})

		t.Run("badArgCount", func(t *testing.T) {
			evalStrToErr(t, `(mapValues (map "a" 1 "b" 2) (map "a" 1 "b" 2))`)
		})
	})

	t.Run("filter", func(t *testing.T) {
		t.Run("basic", func(t *testing.T) {
			// let's make sure the key function is tested here as well
			assertMapValue(
				t,
				evalStrToVal(t, `(mapFilter
					(map "a" 1 "b" 2 "c" 2)
					(fn (k v) (and (strEq k "b") (== v 2))
				  )
				)`),
				map[string]Value{
					"b": &NumberValue{Val: 2},
				},
			)
		})

		t.Run("allowsNils", func(t *testing.T) {
			assertMapValue(
				t,
				evalStrToVal(t, `(mapFilter (map "a" 1 "b" 2) (fn (k v) nil))`),
				map[string]Value{},
			)
		})

		t.Run("badArgCount", func(t *testing.T) {
			evalStrToErr(t, `(mapFilter (map "a" 1 "b" 2))`)
		})

		t.Run("badValue", func(t *testing.T) {
			evalStrToErr(t, `(mapFilter (map "a" 1 "b" nil) (fn (k v) (== v 2)))`)
		})

		t.Run("badReturnValue", func(t *testing.T) {
			evalStrToErr(t, `(mapFilter (map "a" 1 "b" 2) (fn (k v) (+ v 1)))`)
		})

		t.Run("badMapArg", func(t *testing.T) {
			evalStrToErr(t, `(mapFilter "" (fn (k v) (== v 2)))`)
		})

		t.Run("badFnArg", func(t *testing.T) {
			evalStrToErr(t, `(mapFilter (map "a" 1 "b" 2) "")`)
		})
	})

	t.Run("map", func(t *testing.T) {
		t.Run("basic", func(t *testing.T) {
			assertMapValue(
				t,
				evalStrToVal(t, `(mapMap
					(map "a" 1 "b" 2 "c" 2)
					(fn (k v) (if (strEq k "c") (+ v 2) (+ v 1))))`),
				map[string]Value{
					"a": &NumberValue{Val: 2},
					"b": &NumberValue{Val: 3},
					"c": &NumberValue{Val: 4},
				},
			)
		})

		t.Run("badArgCount", func(t *testing.T) {
			evalStrToErr(t, `(mapMap (map "a" 1 "b" 2))`)
		})

		t.Run("badValue", func(t *testing.T) {
			evalStrToErr(t, `(mapMap (map "a" 1 "b" nil) (fn (k v) (+ v 2)))`)
		})

		t.Run("badMapArg", func(t *testing.T) {
			evalStrToErr(t, `(mapMap "" (fn (k v) (== v 2)))`)
		})

		t.Run("badFnArg", func(t *testing.T) {
			evalStrToErr(t, `(mapMap (map "a" 1 "b" 2) "")`)
		})
	})

	t.Run("map", func(t *testing.T) {
		t.Run("basic", func(t *testing.T) {
			assertNumValue(
				t,
				evalStrToVal(t, `(mapReduce 0
					(map "a" 1 "b" 2 "c" 2)
					(fn (t k v) (if (strEq k "c") (+ t (* v 2)) (+ t v ))))`),
				7,
			)
		})

		t.Run("badArgCount", func(t *testing.T) {
			evalStrToErr(t, `(mapReduce 0 (map "a" 1 "b" 2))`)
		})

		t.Run("badValue", func(t *testing.T) {
			evalStrToErr(t, `(mapReduce 0 (map "a" 1 "b" nil) (fn (t k v) (+ t v)))`)
		})

		t.Run("badMapArg", func(t *testing.T) {
			evalStrToErr(t, `(mapReduce 0 "" (fn (t k v) (+ t v)))`)
		})

		t.Run("badFnArg", func(t *testing.T) {
			evalStrToErr(t, `(mapReduce 0 (map "a" 1 "b" 2) "")`)
		})
	})
}

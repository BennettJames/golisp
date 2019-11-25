package golisp2

import (
	"fmt"
	"testing"
)

func Test_string(t *testing.T) {
	type testCase struct {
		name string
		in   string
		out  string
		err  bool
	}

	runCases := func(t *testing.T, cases ...testCase) {
		for i, c := range cases {
			name := c.name
			if len(name) == 0 {
				name = fmt.Sprintf("testCase-%d", i)
			}
			t.Run(name, func(t *testing.T) {
				if c.err {
					evalStrToErr(t, c.in)
				} else {
					assertStringValue(t, evalStrToVal(t, c.in), c.out)
				}
			})
		}
	}

	t.Run("add", func(t *testing.T) {
		runCases(t,
			testCase{
				in:  `(concat "a" "b" "c")`,
				out: "abc",
			},
			testCase{
				in:  `(concat)`,
				out: "",
			},
			testCase{
				in:  `(concat "a" nil)`,
				err: true,
			},
		)
	})
}

func Test_cells(t *testing.T) {

	t.Run("cons", func(t *testing.T) {
		t.Run("basic", func(t *testing.T) {
			assertCellValue(t,
				evalStrToVal(t, `(cons 1 2)`),
				&NumberValue{Val: 1},
				&NumberValue{Val: 2},
			)
		})

		t.Run("nilVals", func(t *testing.T) {
			assertCellValue(t,
				evalStrToVal(t, `(cons)`),
				&NilValue{},
				&NilValue{},
			)
		})

		t.Run("tooManyArgs", func(t *testing.T) {
			evalStrToErr(t, `(cons 1 2 3)`)
		})
	})

	t.Run("car", func(t *testing.T) {
		t.Run("basic", func(t *testing.T) {
			assertNumValue(t,
				evalStrToVal(t, `(car (cons 1 2))`),
				1,
			)
		})

		t.Run("tooManyArgs", func(t *testing.T) {
			evalStrToErr(t, `(car (cons 1 2) (cons 1 2))`)
		})

		t.Run("badType", func(t *testing.T) {
			evalStrToErr(t, `(car "abc")`)
		})
	})

	t.Run("cdr", func(t *testing.T) {
		t.Run("basic", func(t *testing.T) {
			assertNumValue(t,
				evalStrToVal(t, `(cdr (cons 1 2))`),
				2,
			)
		})

		t.Run("tooManyArgs", func(t *testing.T) {
			evalStrToErr(t, `(cdr (cons 1 2) (cons 1 2))`)
		})

		t.Run("badType", func(t *testing.T) {
			evalStrToErr(t, `(cdr "abc")`)
		})
	})
}

func Test_math(t *testing.T) {
	type testCase struct {
		name string
		in   string
		out  float64
		err  bool
	}

	runCases := func(t *testing.T, cases ...testCase) {
		for i, c := range cases {
			name := c.name
			if len(name) == 0 {
				name = fmt.Sprintf("testCase-%d", i)
			}
			t.Run(name, func(t *testing.T) {
				if c.err {
					evalStrToErr(t, c.in)
				} else {
					assertNumValue(t, evalStrToVal(t, c.in), c.out)
				}
			})
		}
	}

	t.Run("add", func(t *testing.T) {
		runCases(t,
			testCase{
				in:  `(+ 6 1.5)`,
				out: 7.5,
			},
			testCase{
				in:  `(+ 5 3 -2)`,
				out: 6,
			},
			testCase{
				in:  `(+ 1 nil)`,
				err: true,
			},
			testCase{
				in:  `(+)`,
				err: true,
			},
		)
	})

	t.Run("sub", func(t *testing.T) {
		runCases(t,
			testCase{
				in:  `(- 6 1.5)`,
				out: 4.5,
			},
			testCase{
				in:  `(- 5 3 3)`,
				out: -1,
			},
			testCase{
				in:  `(- 22)`,
				out: -22,
			},
			testCase{
				in:  `(- 1 nil)`,
				err: true,
			},
			testCase{
				in:  `(-)`,
				err: true,
			},
		)
	})

	t.Run("mult", func(t *testing.T) {
		runCases(t,
			testCase{
				in:  `(* 6 1.5)`,
				out: 9,
			},
			testCase{
				in:  `(* 5 2 2)`,
				out: 20,
			},
			testCase{
				in:  `(* 1 nil)`,
				err: true,
			},
		)
	})

	t.Run("div", func(t *testing.T) {
		runCases(t,
			testCase{
				in:  `(/ 6 2)`,
				out: 3,
			},
			testCase{
				in:  `(/ 5 2 2)`,
				out: 1.25,
			},
			testCase{
				in:  `(/ 1 nil)`,
				err: true,
			},
		)
	})
}

func Test_comparisons(t *testing.T) {
	type testCase struct {
		name string
		in   string
		out  bool
		err  bool
	}

	runCases := func(t *testing.T, cases ...testCase) {
		for i, c := range cases {
			name := c.name
			if len(name) == 0 {
				name = fmt.Sprintf("testCase-%d", i)
			}
			t.Run(name, func(t *testing.T) {
				if c.err {
					evalStrToErr(t, c.in)
				} else {
					assertBoolValue(t, evalStrToVal(t, c.in), c.out)
				}
			})
		}
	}

	t.Run("and", func(t *testing.T) {
		runCases(t,
			testCase{
				in:  `(and false)`,
				out: false,
			},
			testCase{
				in:  `(and true true true)`,
				out: true,
			},
			testCase{
				in:  `(and true true true false)`,
				out: false,
			},
			testCase{
				in:  `(and true "abc")`,
				err: true,
			},
			testCase{
				in:  `(and)`,
				err: true,
			},
		)
	})

	t.Run("or", func(t *testing.T) {
		runCases(t,
			testCase{
				in:  `(or false)`,
				out: false,
			},
			testCase{
				in:  `(or false false true)`,
				out: true,
			},
			testCase{
				in:  `(or true "abc")`,
				err: true,
			},
			testCase{
				in:  `(or)`,
				err: true,
			},
		)
	})

	t.Run("not", func(t *testing.T) {
		runCases(t,
			testCase{
				in:  `(not true)`,
				out: false,
			},
			testCase{
				in:  `(not false)`,
				out: true,
			},
			testCase{
				in:  `(not "abc")`,
				err: true,
			},
			testCase{
				in:  `(not)`,
				err: true,
			},
			testCase{
				in:  `(not false false)`,
				err: true,
			},
		)
	})

	t.Run("eq", func(t *testing.T) {
		runCases(t,
			testCase{
				in:  `(== 1 2)`,
				out: false,
			},
			testCase{
				in:  `(== 1 1)`,
				out: true,
			},
			testCase{
				in:  `(== 1 nil)`,
				err: true,
			},
		)
	})

	t.Run("strEq", func(t *testing.T) {
		runCases(t,
			testCase{
				in:  `(strEq "a" "b")`,
				out: false,
			},
			testCase{
				in:  `(strEq "a" "a")`,
				out: true,
			},
			testCase{
				in:  `(strEq "a" nil)`,
				err: true,
			},
			testCase{
				in:  `(strEq "a")`,
				err: true,
			},
			testCase{
				in:  `(strEq "a" "b" "c")`,
				err: true,
			},
		)
	})

	t.Run("gt", func(t *testing.T) {
		runCases(t,
			testCase{
				in:  `(> 1 2)`,
				out: false,
			},
			testCase{
				in:  `(> 1 0)`,
				out: true,
			},
			testCase{
				in:  `(> 1 1)`,
				out: false,
			},
			testCase{
				in:  `(> 1 nil)`,
				err: true,
			},
		)
	})

	t.Run("lt", func(t *testing.T) {
		runCases(t,
			testCase{
				in:  `(< 1 2)`,
				out: true,
			},
			testCase{
				in:  `(< 1 0)`,
				out: false,
			},
			testCase{
				in:  `(< 1 1)`,
				out: false,
			},
			testCase{
				in:  `(< 1 nil)`,
				err: true,
			},
		)
	})

	t.Run("gte", func(t *testing.T) {
		runCases(t,
			testCase{
				in:  `(>= 1 2)`,
				out: false,
			},
			testCase{
				in:  `(>= 1 0)`,
				out: true,
			},
			testCase{
				in:  `(>= 1 1)`,
				out: true,
			},
			testCase{
				in:  `(>= 1 nil)`,
				err: true,
			},
		)
	})

	t.Run("lte", func(t *testing.T) {
		runCases(t,
			testCase{
				in:  `(<= 1 2)`,
				out: true,
			},
			testCase{
				in:  `(<= 1 0)`,
				out: false,
			},
			testCase{
				in:  `(<= 1 1)`,
				out: true,
			},
			testCase{
				in:  `(<= 1 nil)`,
				err: true,
			},
		)
	})
}

func Test_print(t *testing.T) {
	// note (bs): this isn't really a meaningful test; not sure if there's a good
	// way to do so without some very awkward dependency reconfiguration

	assertNilValue(t, evalStrToVal(t, `(print (list 1 2 3))`))
	assertNilValue(t, evalStrToVal(t, `(print)`))
	assertNilValue(t, evalStrToVal(t, `(print 1 2 3)`))
}

func Test_len(t *testing.T) {

	t.Run("list", func(t *testing.T) {
		assertNumValue(t, evalStrToVal(t, `(len (list 1 2 3))`), 3)
	})

	t.Run("map", func(t *testing.T) {
		assertNumValue(t, evalStrToVal(t, `(len (map "a" 1 "b" 2))`), 2)
	})

	t.Run("string", func(t *testing.T) {
		assertNumValue(t, evalStrToVal(t, `(len "abcde")`), 5)
	})

	t.Run("badType", func(t *testing.T) {
		evalStrToErr(t, `(len nil)`)
	})

	t.Run("badArgLen", func(t *testing.T) {
		evalStrToErr(t, `(len "a" "b")`)
	})
}

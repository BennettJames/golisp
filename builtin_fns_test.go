package golisp2

import (
	"fmt"
	"testing"
)

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

	t.Run("div", func(t *testing.T) {
		runCases(t,
			testCase{
				in:  `(/ 6 2)`,
				out: 3,
			},
			testCase{
				in:  `(/ 5 2)`,
				out: 2.5,
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

package golisp2

import "testing"

func Test_ExecString(t *testing.T) {
	testCases := []struct {
		Name   string
		Input  string
		Output string
		Err    error
	}{
		{
			Name:   "SuperSimple",
			Input:  `(+ 1 2)`,
			Output: `3`,
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			resStr, resErr := ExecString(c.Input)
			if c.Err != nil && resErr == nil {
				t.Fatalf("Expected error [%s], got none [output=%s]", c.Err, resStr)
			}
			if c.Err == nil && resErr != nil {
				t.Fatalf("Expected no error [output=%s], got none [error=%s]", c.Output, resErr)
			}
			if resStr != c.Output {
				// note (bs): may want some fuzz-toleration here for things like whitespace
				t.Fatalf("Unexpected output [expected=%s] [actual=%s]", c.Output, resStr)
			}
		})
	}
}

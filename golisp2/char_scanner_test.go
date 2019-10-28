package golisp2

import (
	"testing"
)

func Test_CharScanner(t *testing.T) {
	// honestly I don't think this needs a test table. It's pretty dumb given that
	// it doesn't even use a proper reader under the hood.

	scanner := NewCharScanner("hi \x00 😊")
	expected := []rune{
		'h',
		'i',
		' ',
		'\x00',
		' ',
		'😊',
	}
	for i, expectedR := range expected {
		peekR, peekOk := scanner.Peek()
		nextR, nextOk := scanner.Next()
		if !peekOk || !nextOk {
			t.Fatal("unexpected end to scan")
		}
		if expectedR != peekR {
			t.Fatalf("peek returned wrong value [index=%d] [expected=%s] [actual=%s]",
				i, string(expectedR), string(peekR))
		}
		if expectedR != nextR {
			t.Fatalf("next returned wrong value [index=%d] [expected=%s] [actual=%s]",
				i, string(expectedR), string(nextR))
		}
	}
	if !scanner.Done() {
		t.Fatal("scanner did not complete")
	}
	if _, finalOk := scanner.Next(); finalOk {
		t.Fatal("next should return false after completion")
	}
}

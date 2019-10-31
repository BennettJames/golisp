package golisp2

import (
	"strings"
	"testing"
)

func Test_RuneScanner(t *testing.T) {
	// honestly I don't think this needs a test table. It's pretty dumb given that
	// it doesn't even use a proper reader under the hood.

	rs := NewRuneScanner(strings.NewReader("hi \x00 😊"))
	expected := []rune{
		'h',
		'i',
		' ',
		'\x00',
		' ',
		'😊',
	}
	for i, expectedR := range expected {
		rs.Advance()
		if rs.Done() {
			t.Fatal("unexpected end to scan")
		}
		if expectedR != rs.Rune() {
			t.Fatalf("scanner returned wrong value [index=%d] [expected=%s] [actual=%s]",
				i, string(expectedR), string(rs.Rune()))
		}
	}
	rs.Advance()
	if !rs.Done() {
		t.Fatal("scanner did not complete")
	}
}

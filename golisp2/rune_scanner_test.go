package golisp2

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_RuneScanner(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		rs := NewRuneScanner("testfile", strings.NewReader("hi \x00 😊"))
		expected := []rune{
			'h',
			'i',
			' ',
			'\x00',
			' ',
			'😊',
		}
		for _, expectedR := range expected {
			rs.Advance()
			require.False(t, rs.Done(), "unexpected end to scan")
			require.Equal(t, expectedR, rs.Rune())
		}
		rs.Advance()
		require.True(t, rs.Done(), "scanned should complete after everyting's read")
		require.Equal(t, io.EOF, rs.Err())
	})

	t.Run("pos", func(t *testing.T) {
		fName := "itsATestFile.l"
		rs := NewRuneScanner(fName, strings.NewReader("(\n+ 1 2\n)"))
		expected := []struct {
			r        rune
			col, row int
		}{
			{'(', 1, 1},
			{'\n', 2, 1},
			{'+', 1, 2},
			{' ', 2, 2},
			{'1', 3, 2},
			{' ', 4, 2},
			{'2', 5, 2},
			{'\n', 6, 2},
			{')', 1, 3},
		}
		for _, e := range expected {
			rs.Advance()
			require.False(t, rs.Done(), "unexpected end to scan")
			require.Equal(t, e.r, rs.Rune())
			require.Equal(t, fName, rs.Pos().SourceFile)
			require.Equal(t, e.col, rs.Pos().Col)
			require.Equal(t, e.row, rs.Pos().Row)
		}
	})
}

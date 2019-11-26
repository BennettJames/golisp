package golisp2

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_RuneScanner(t *testing.T) {
	fName := "itsATestFile.l"
	t.Run("basicScan", func(t *testing.T) {
		rs := NewRuneScanner(fName, strings.NewReader("(\n+ 1 \t😊\n)"))
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
			{'\t', 5, 2},
			{'😊', 6, 2},
			{'\n', 7, 2},
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
		rs.Advance()
		require.True(t, rs.Done(), "scanned should complete after everyting's read")
		require.Equal(t, io.EOF, rs.Err())
	})

	t.Run("forbiddenChar", func(t *testing.T) {
		rs := NewRuneScanner(fName, strings.NewReader("\x00abc"))
		rs.Advance()
		require.Error(t, rs.Err())
		asForbidden, isForbidden := rs.Err().(*ForbiddenRuneError)
		require.True(t, isForbidden)
		require.Equal(t, '\x00', asForbidden.R)
		require.Equal(t, ScannerPosition{
			SourceFile: fName,
			Col:        1,
			Row:        1,
		}, asForbidden.Pos)
	})
}

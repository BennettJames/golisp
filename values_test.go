package golisp2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_boolValue(t *testing.T) {
	t.Run("boolInspectStr", func(t *testing.T) {
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

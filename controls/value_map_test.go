package controls

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewMap(t *testing.T) {
	m, err := NewMap("tilt:0", "pan:64")
	require.NoError(t, err)

	require.Equal(t, byte(0), m["tilt"])
	require.Equal(t, byte(64), m["pan"])
}

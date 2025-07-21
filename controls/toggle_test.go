package controls

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewToggle(t *testing.T) {
	toggle := NewToggle()
	toggle.On()
	require.Equal(t, toggle.GetValue(), true)

	var control Control = toggle
	err := control.SetValueString("true")
	require.NoError(t, err)
	require.Equal(t, toggle.GetValue(), true)
}

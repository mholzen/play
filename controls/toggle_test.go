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
	control.SetValueString("true")
	require.Equal(t, toggle.GetValue(), true)
}

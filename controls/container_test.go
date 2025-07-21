package controls

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ToggleInContainer(t *testing.T) {
	toggle := NewToggle()
	list := NewList(1)
	list.SetItem(0, toggle)

	// Get item back as Control interface
	item, err := list.GetItem("0")
	require.NoError(t, err)

	control, ok := item.(Control)
	require.True(t, ok, "Item should implement Control interface")

	// Verify control methods work
	err = control.SetValueString("true")
	require.NoError(t, err)
	require.Equal(t, true, toggle.GetValue())
	require.Equal(t, "true", control.GetValueString())
}

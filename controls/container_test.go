package controls

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ToggleInContainer(t *testing.T) {
	// Create a toggle and container
	toggle := NewToggle()
	list := NewList(1)
	list.SetItem(0, toggle)

	// Get item back as Control interface
	item, err := list.GetItem("0")
	require.NoError(t, err)

	control, ok := item.(Control)
	require.True(t, ok, "Item should implement Control interface")

	// Verify control methods work
	control.SetValueString("true")
	require.Equal(t, true, toggle.GetValue())
	require.Equal(t, "true", control.GetValueString())
}

// func Test_MuxInContainer(t *testing.T) {
// 	// Create a mux and container
// 	mux := NewMux[bool]()
// 	toggle1 := NewToggle()
// 	toggle2 := NewToggle()
// 	mux.Add("toggle1", toggle1)
// 	mux.Add("toggle2", toggle2)
// 	list := NewList(1)
// 	list.SetItem(0, &mux)

// 	// Get item back as Control interface
// 	item, err := list.GetItem("0")
// 	require.NoError(t, err)

// 	control, ok := item.(Control)
// 	require.True(t, ok, "Item should implement Control interface")

// 	// Verify control methods work
// 	control.SetValueString("toggle2")
// 	require.Equal(t, "toggle2", control.GetValueString())
// 	require.Equal(t, "toggle2", mux.GetSource())
// }

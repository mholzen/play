package controls

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewSelector(t *testing.T) {
	selector := NewSelector[string]()
	require.NotNil(t, selector)
	require.NotNil(t, selector.Options)
	require.Equal(t, "", selector.Selected)
}

func Test_SetOptions(t *testing.T) {
	selector := NewSelector[string]()
	options := map[string]string{
		"option1": "value1",
		"option2": "value2",
	}

	selector.SetOptions(options)
	require.Equal(t, options, selector.Options)
}

func Test_SetSelected_Valid(t *testing.T) {
	selector := NewSelector[string]()
	options := map[string]string{
		"option1": "value1",
		"option2": "value2",
	}
	selector.SetOptions(options)

	err := selector.SetSelected("option1")
	require.NoError(t, err)
	require.Equal(t, "option1", selector.GetSelected())
}

func Test_SetSelected_Invalid(t *testing.T) {
	selector := NewSelector[string]()
	options := map[string]string{
		"option1": "value1",
	}
	selector.SetOptions(options)

	require.Panics(t, func() {
		selector.SetSelected("invalid")
	})
}

func Test_GetSelectedValue(t *testing.T) {
	selector := NewSelector[int]()
	options := map[string]int{
		"first":  10,
		"second": 20,
	}
	selector.SetOptions(options)
	selector.SetSelected("second")

	require.Equal(t, 20, selector.GetSelectedValue())
}

func Test_MarshalJSON(t *testing.T) {
	selector := NewSelector[string]()
	options := map[string]string{
		"b_option": "value2",
		"a_option": "value1",
	}
	selector.SetOptions(options)
	selector.SetSelected("a_option")

	data, err := selector.MarshalJSON()
	require.NoError(t, err)
	require.Contains(t, string(data), `"value":"a_option"`)
	require.Contains(t, string(data), `"options":["a_option","b_option"]`)
}

func Test_GetValueString(t *testing.T) {
	selector := NewSelector[string]()
	options := map[string]string{
		"test": "value",
	}
	selector.SetOptions(options)
	selector.SetSelected("test")

	require.Equal(t, "test", selector.GetValueString())
}

func Test_SetValueString(t *testing.T) {
	selector := NewSelector[string]()
	options := map[string]string{
		"test": "value",
	}
	selector.SetOptions(options)

	notified := false
	var notifiedValue string
	OnChange(selector, func(value string) {
		notified = true
		notifiedValue = value
	})

	err := selector.SetValueString("test")
	require.NoError(t, err)
	require.Equal(t, "test", selector.GetSelected())
	require.True(t, notified)
	require.Equal(t, "value", notifiedValue)
}

package controls

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNumericDial(t *testing.T) {
	d := NewNumericDial()
	d.SetValue(100)
	assert.Equal(t, byte(100), d.Value)

	var _ Control = d
}

func Test_NumericDial(t *testing.T) {
	numericDial := NewNumericDial()

	var _ Control = numericDial
	var _ Settable = numericDial
	var _ Dial[byte] = numericDial
}

func Test_ObservableNumericDial(t *testing.T) {
	observableNumericDial := NewObservableNumericalDial()

	var _ Control = observableNumericDial

	ch := make(chan byte)
	observableNumericDial.AddObserver(ch)

	go func() {
		observableNumericDial.SetValue(42)
	}()

	value := <-ch
	require.Equal(t, byte(42), value)
}

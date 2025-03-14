package controls

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NumericDial(t *testing.T) {
	numericDial := NewNumericDial()

	var _ Control = numericDial
	var _ Settable = numericDial
	var _ Dial[byte] = numericDial
}

func Test_ObservableNumericDial(t *testing.T) {
	observableNumericDial := NewObservableNumericalDial()

	var _ Control = observableNumericDial

	testValidated := make(chan int)
	ch := make(chan byte)
	observableNumericDial.AddObserver(ch)

	go func() {
		observableNumericDial.SetValue(42)
		testValidated <- 1
	}()

	value := <-ch
	require.Equal(t, byte(42), value)
	<-testValidated
}

package controls

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Mux(t *testing.T) {
	in1 := NewToggle()
	in2 := NewToggle()
	in2.On()

	mux := NewMux[bool]()
	mux.Add("in1", in1)
	mux.Add("in2", in2)

	mux.SetSource("in2")
	require.Equal(t, mux.GeValue(), true)

	mux.SetSource("in1")
	require.Equal(t, mux.GeValue(), false)
}

func Test_Mux_ValueMap(t *testing.T) {
	inA := NewNumericDialMap("ch1", "ch2")
	inB := NewNumericDialMap("ch1", "ch2")
	inB.SetValue(ValueMap{"ch1": 0, "ch2": 0})
	inB.SetValue(ValueMap{"ch1": 255, "ch2": 255})

	mux := NewMux[ValueMap]()
	mux.Add("inA", inA)
	mux.Add("inB", inB)

	mux.SetSource("inA")
	require.Equal(t, mux.GeValue()["ch1"], uint8(0))

	mux.SetSource("inB")
	require.Equal(t, mux.GeValue()["ch1"], uint8(255))

	muxChannel := mux.Channel()
	done := make(chan bool)
	go func() {
		val := <-muxChannel
		require.Equal(t, val["ch1"], uint8(2))
		done <- true

		val = <-muxChannel
		require.Equal(t, val["ch1"], uint8(2))
		done <- true
	}()

	inA.SetValue(ValueMap{"ch1": 1, "ch2": 1}) // should not emit
	inB.SetValue(ValueMap{"ch1": 2, "ch2": 2})
	<-done

	mux.SetSource("inA")
	inB.SetValue(ValueMap{"ch1": 2, "ch2": 2}) // should not emit
	inA.SetValue(ValueMap{"ch1": 1, "ch2": 1})
	<-done
}

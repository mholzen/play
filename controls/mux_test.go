package controls

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// func Test_Mux(t *testing.T) {
// 	in1 := NewToggle()
// 	in2 := NewToggle()
// 	in2.On()

// 	mux := NewMux[bool]()
// 	mux.Add("in1", in1)
// 	mux.Add("in2", in2)

// 	mux.SetSource("in2")
// 	require.Equal(t, mux.GetValue(), true)

// 	mux.SetSource("in1")
// 	require.Equal(t, mux.GetValue(), false)
// }

func Test_Mux_ValueMap(t *testing.T) {
	inA := NewObservableNumericDialMap("ch1")
	inB := NewObservableNumericDialMap("ch1")

	mux := NewMux[ChannelValues]()
	mux.Add("inA", inA)
	mux.Add("inB", inB)

	mux.SetSource("inA")
	muxChannel := make(chan ChannelValues)
	mux.AddObserver(muxChannel)

	done := make(chan bool)
	go func() {
		val := <-muxChannel
		assert.Equal(t, val["ch1"], uint8(42))
		<-done

		val = <-muxChannel
		assert.Equal(t, val["ch1"], uint8(43))
		<-done

	}()

	inA.SetValue(ChannelValues{"ch1": 42})
	done <- true

	mux.SetSource("inB") // Changing source does not trigger notification
	inB.SetValue(ChannelValues{"ch1": 43})
	done <- true
}

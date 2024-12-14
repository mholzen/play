package controls

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

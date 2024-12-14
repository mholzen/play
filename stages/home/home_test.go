package home

import (
	"log"
	"testing"

	"github.com/mholzen/play-go/controls"
	"github.com/mholzen/play-go/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_RootSurfaceHome(t *testing.T) {

	clock := controls.NewClock(120)
	clock.Start()
	rootSurface := GetRootSurface(home.Universe, clock)

	// dials
	item, err := rootSurface.GetItem("0")
	require.NoError(t, err)
	dials, ok := item.(*controls.ObservableDialMap)
	require.True(t, ok)
	assert.Contains(t, dials.GetValue(), "r")

	// mux
	item, err = rootSurface.GetItem("2")
	require.NoError(t, err)
	mux, ok := item.(*controls.Mux[fixture.FixtureValues])
	require.True(t, ok)
	require.NoError(t, mux.SetSource("dials"))

	advance := make(chan int)

	muxChannel := make(chan fixture.FixtureValues)
	mux.AddObserver(muxChannel)

	go func() {
		dials.SetChannelValue("r", 42)
		<-advance // 1

		mux.SetSource("rainbow")
		dials.SetChannelValue("r", 43) // should not have an effect
		// rainbow is not setting values through the mux channel
		log.Printf("set r to 43")
		<-advance // 2
	}()

	value := <-muxChannel

	address := home.TomeShine.GetAddresses()[0]
	expected := value[address]["r"]
	assert.Equal(t, byte(42), expected)

	advance <- 1

	// wait for rainbow to set a non zero value
	for i := 0; i < 100; i++ {
		value = <-muxChannel
		expected = value[address]["tilt"]
		if expected > 0 {
			break
		}
	}
	assert.Greater(t, expected, byte(0))
}

package home

import (
	"log"
	"testing"

	"github.com/mholzen/play-go/controls"
	"github.com/mholzen/play-go/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_RootSurfaceDialControls(t *testing.T) {
	clock := controls.NewClock(120)
	clock.Start()
	rootSurface := GetRootSurface(home.Universe, clock)

	// dial map
	item, err := rootSurface.GetItem("1")
	require.NoError(t, err)

	container, ok := item.(controls.Container)
	require.True(t, ok)
	item, err = container.GetItem("r")
	require.NoError(t, err)
	assert.NotNil(t, item)

	// dial list
	item, err = rootSurface.GetItem("1")
	require.NoError(t, err)
	dialList, ok := item.(controls.OrderedContainer)
	require.True(t, ok)
	assert.Contains(t, dialList.Keys(), "r")
}

func Test_RootSurfaceMux(t *testing.T) {
	clock := controls.NewClock(120)
	clock.Start()
	rootSurface := GetRootSurface(home.Universe, clock)

	// Get dialMap
	item, err := rootSurface.GetItem("1")
	require.NoError(t, err)

	require.IsType(t, &controls.DialList{}, item)

	dialList, ok := item.(*controls.DialList)
	require.True(t, ok)

	dial, err := dialList.GetItem("r")
	require.NoError(t, err)

	redDial, ok := dial.(*controls.ObservableNumericalDial)
	require.True(t, ok)

	// mux
	item, err = rootSurface.GetItem("0")
	require.NoError(t, err)
	mux, ok := item.(*controls.Mux[fixture.FixtureValues])
	require.True(t, ok)
	require.NoError(t, mux.SetSource("dials"))

	advance := make(chan int)

	muxChannel := make(chan fixture.FixtureValues)
	mux.AddObserver(muxChannel)

	go func() {
		t.Helper() // Marks this function as a test helper

		dialList.SetChannelValue("r", 0xa)
		<-advance // wait for 1

		redDial.SetValue(0xb)
		log.Printf("set r to 0xb")
		<-advance // wait for 2

		mux.SetSource("rainbow")
		dialList.SetChannelValue("r", 0xc) // should not have an effect
		// rainbow is not setting values through the mux channel
		log.Printf("set r to 0xc")
		<-advance // wait for 3
	}()

	value := <-muxChannel

	address := home.TomeShine.GetAddresses()[0]
	expected := value[address]["r"]
	assert.Equal(t, byte(10), expected)

	advance <- 1

	value = <-muxChannel

	expected = value[address]["r"]
	assert.Equal(t, byte(0xb), expected)

	advance <- 2

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

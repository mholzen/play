package home

import (
	"testing"

	"github.com/mholzen/play-go/controls"
	"github.com/mholzen/play-go/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_RootSurfaceHome(t *testing.T) {

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

	advance := make(chan bool)

	muxChannel := make(chan fixture.FixtureValues)
	mux.AddObserver(muxChannel)

	go func() {
		dials.SetChannelValue("r", 42)
		<-advance
	}()

	value := <-muxChannel

	address := home.TomeShine[1].Address
	expected := value[address]["r"]
	assert.Equal(t, byte(42), expected)

	advance <- true

	// require.NoError(t, mux.SetSource("rainbow"))
	// advance <- true
}

package home

import (
	"log"
	"testing"

	"github.com/mholzen/play-go/controls"
	"github.com/mholzen/play-go/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func _Test_RootSurfaceHome(t *testing.T) {

	rootSurface := GetRootSurface(home.Universe, clock)

	item, err := rootSurface.GetItem("0")
	require.NoError(t, err)
	dials, ok := item.(*controls.DialMap)
	require.True(t, ok)

	// the dialmap messages are sent to the fixturevalues, not to the mux
	// need multiple listeners

	item, err = rootSurface.GetItem("2")
	require.NoError(t, err)
	mux, ok := item.(*controls.Mux[fixture.FixtureValues])
	require.True(t, ok)
	require.NoError(t, mux.SetSource("dials"))

	move := make(chan bool)
	received := false
	go func() {
		<-move
		log.Printf("waiting for value")
		value := <-mux.Channel()
		expected := value[home.TomeShine[0].Address]["r"]
		assert.Equal(t, 42, expected)
		received = true
		move <- true
	}()

	move <- true
	for !received {
		dials.SetChannelValue("r", 42)
	}
}

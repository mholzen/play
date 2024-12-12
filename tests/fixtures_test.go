package fixture

import (
	"testing"

	"github.com/mholzen/play-go/fixture"
	"github.com/stretchr/testify/assert"
)

func Test_Fixtures_SetChannelValue(t *testing.T) {
	f := fixture.NewFixtures()
	f.AddFixture(fixture.NewFreedomPar(), 1)
	f.AddFixture(fixture.NewFreedomPar(), 100)

	f.SetChannelValue("r", 1)

	assert.Equal(t, byte(1), f[1].Fixture.GetValueMap()["r"])
	assert.Equal(t, byte(1), f[100].Fixture.GetValueMap()["r"])
}

// func Test_Fixtures_IsEmitter(t *testing.T) {
// 	f := fixture.NewFixtures()

// 	var _ controls.Emitter[fixture.FixtureValues] = f
// }

// func Test_Mux_FixtureValues(t *testing.T) {
// 	// Create two fixture emitters
// 	fixtures1 := fixture.NewFixtures()
// 	fixtures2 := fixture.NewFixtures()

// 	f1 := fixture.NewFreedomPar()
// 	f2 := fixture.NewFreedomPar()
// 	fixtures1.AddFixture(f1, 1)
// 	fixtures2.AddFixture(f2, 1)

// 	emitter1 := fixture.NewFixtureEmitter(fixtures1)
// 	emitter2 := fixture.NewFixtureEmitter(fixtures2)

// 	// Create and configure mux
// 	mux := controls.NewMux[fixture.FixtureValues]()
// 	mux.Add("em1", &emitter1)
// 	mux.Add("em2", &emitter2)

// 	// Test initial values
// 	mux.SetSource("em1")
// 	require.Equal(t, mux.GetSource(), "em1")

// 	// Test value changes propagate through mux
// 	done := make(chan bool)
// 	muxChan := mux.Channel()

// 	go func() {
// 		val := <-muxChan
// 		require.Equal(t, val[1]["r"], byte(100))
// 		done <- true

// 		val = <-muxChan
// 		require.Equal(t, val[1]["r"], byte(200))
// 		done <- true
// 	}()

// 	emitter2.SetChannelValue("r", 50) // Should not emit since em1 is selected
// 	emitter1.SetChannelValue("r", 100)
// 	<-done

// 	mux.SetSource("em2")
// 	emitter1.SetChannelValue("r", 150) // Should not emit since em2 is selected
// 	emitter2.SetChannelValue("r", 200)
// 	<-done

// }

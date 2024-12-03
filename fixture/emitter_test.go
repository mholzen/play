package fixture

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ChannelIsSaved(t *testing.T) {
	fixture := NewFreedomPar()
	fixtures := NewFixtures()
	fixtures.AddFixture(fixture, 1)
	emitter := NewFixtureEmitter(fixtures)

	done := make(chan bool)
	ch1 := emitter.Channel()
	go func() {
		done <- true
		val := <-ch1
		require.Equal(t, val[1]["r"], byte(1))
		done <- true
	}()

	<-done
	emitter.SetChannelValue("r", 1)
	<-done
}

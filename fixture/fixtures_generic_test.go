package fixture

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FixtureGeneric_SetChannelValue(t *testing.T) {
	f := NewFixturesGeneric[Fixture]()
	f.AddFixture(NewFreedomPar(), 1)
	f.AddFixture(NewFreedomPar(), 100)

	f.SetChannelValue("r", 1)

	f1 := f.GetFixtures()[1]
	f100 := f.GetFixtures()[100]

	assert.Equal(t, byte(1), f1.GetValueMap()["r"])
	assert.Equal(t, byte(1), f100.GetValueMap()["r"])

	values := f.GetByteArray()
	assert.Equal(t, byte(1), values[1+1]) // red is channel 1 (dimmer is channel 0)
	assert.Equal(t, byte(1), values[100+1])
}

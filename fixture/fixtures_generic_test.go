package fixture

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FixtureGeneric_SetChannelValue(t *testing.T) {
	f := NewAddressableFixtures[ChannelFixture]()
	f.AddFixture(NewFreedomPar(), 1)
	f.AddFixture(NewFreedomPar(), 100)

	f.SetChannelValue("r", 1)

	f1 := f.GetFixtures()[1]
	f100 := f.GetFixtures()[100]

	assert.Equal(t, byte(1), f1.GetChannelValues()["r"])
	assert.Equal(t, byte(1), f100.GetChannelValues()["r"])

	values := f.GetByteArray()
	assert.Equal(t, byte(1), values[1+1]) // red is channel 1 (dimmer is channel 0)
	assert.Equal(t, byte(1), values[100+1])
}

func Test_Fixtures_Separate(t *testing.T) {
	fs1 := *NewAddressableFixtures[ChannelFixture]()
	fs1.AddFixture(NewFreedomPar(), 1)
	assert.Equal(t, byte(0), fs1.GetChannelValues()["r"])

	fs2 := fs1.Clone()

	f1 := fs1.GetFixtures()[1]
	f2 := fs2.GetFixtures()[1]

	f1.SetChannelValue("r", 1)

	assert.Equal(t, byte(1), f1.GetChannelValues()["r"])
	assert.Equal(t, byte(0), f2.GetChannelValues()["r"]) // f2 should not be affected by f1
}

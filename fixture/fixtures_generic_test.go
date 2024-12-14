package fixture

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FixtureGeneric_SetChannelValue(t *testing.T) {
	f := NewFixturesGeneric[ChannelFixture]()
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
	fs1 := *NewFixturesGeneric[ChannelFixture]()
	fs1.AddFixture(NewFreedomPar(), 1)
	assert.Equal(t, byte(0), fs1.GetChannelValues()["r"])

	fs2 := fs1.Clone()

	f1 := fs1.GetFixtures()[1]
	f2 := fs2.GetFixtures()[1]

	f1.SetChannelValue("r", 1)

	assert.Equal(t, byte(1), f1.GetChannelValues()["r"])
	assert.Equal(t, byte(0), f2.GetChannelValues()["r"]) // f2 should not be affected by f1
}

type Stringer interface {
	Foo() string
}

type StringerImpl []int

func (s StringerImpl) Foo() string {
	return fmt.Sprintf("%v", s)
}

func Test_Basic(t *testing.T) {
	concrete1 := StringerImpl{1}
	var s1 Stringer = concrete1
	assert.Equal(t, "[1]", s1.Foo())

	var s2 Stringer = s1
	assert.Equal(t, "[1]", s2.Foo())
	s2 = StringerImpl{2}
	assert.Equal(t, "[2]", s2.Foo())

	assert.Equal(t, "[1]", s1.Foo())
}

func Test_Basic_List(t *testing.T) {
	list1 := make(map[int]int)
	list1[1] = 1
	list1[2] = 2

	list2 := list1

	list2[1] = 3

	assert.Equal(t, list1, list2)
	assert.Equal(t, list1[1], 3)
}

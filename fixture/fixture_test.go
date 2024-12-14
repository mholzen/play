package fixture

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Fixtures_SetChannelValue(t *testing.T) {
	f := NewFreedomPar()

	var fixture Fixture = &f

	fixture.SetChannelValue("r", 1)

	assert.Equal(t, byte(1), fixture.GetChannelValues()["r"])
}

func Test_Go_Map_Increment(t *testing.T) {
	a := make(map[int]int)
	a[1] = 1
	a[1]++
	assert.Equal(t, a[1], 2)
}

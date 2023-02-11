package fixture

import (
	"github.com/akualab/dmx"
)

func NewFixtureList() FixtureList {
	return FixtureList{
		Fixtures:  make([]FixtureI, 0),
		Addresses: make([]int, 0),
	}
}

type FixtureList struct {
	Fixtures  []FixtureI
	Addresses []int
}

func (f *FixtureList) AddFixture(fixture FixtureI, address int) {
	f.Fixtures = append(f.Fixtures, fixture)
	f.Addresses = append(f.Addresses, address)
	// TODO: check for overlap
}

func (f FixtureList) SetValue(name string, value byte) {
	for _, fixture := range f.Fixtures {
		fixture.SetValue(name, value)
	}
}

func (f FixtureList) SetAll(value byte) {
	for _, fixture := range f.Fixtures {
		fixture.SetAll(value)
	}
}

func (f FixtureList) GetLastAddress() int {
	lastIndex := len(f.Addresses) - 1
	lastStartAddress := f.Addresses[lastIndex]
	lastSize := len(f.Fixtures[lastIndex].GetValues())
	return lastStartAddress + lastSize
}

func (f FixtureList) GetValues() []byte {
	res := make([]byte, f.GetLastAddress())
	for i, fixture := range f.Fixtures {
		values := fixture.GetValues()
		for j, value := range values {
			res[f.Addresses[i]+j] = value
		}
	}
	return res
}

func (f *FixtureList) Render(connection dmx.DMX) error {
	for i, fixture := range f.Fixtures {
		for j, value := range fixture.GetValues() {
			channel := f.Addresses[i] + j
			connection.SetChannel(channel, value)
			// log.Printf("set channel %d to value %d", channel, value)
		}
	}
	// log.Print("rendering")
	return connection.Render()
}

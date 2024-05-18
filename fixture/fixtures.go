package fixture

import (
	"time"

	"github.com/akualab/dmx"
)

func NewFixtures() Fixtures {
	return Fixtures{
		Fixtures:  make([]FixtureI, 0),
		Addresses: make([]int, 0),
	}
}
func NewFixtureList(fixture FixtureI, address ...int) Fixtures {
	res := NewFixtures()
	for _, a := range address {
		res.AddFixture(fixture, a)
	}
	return res
}

type Fixtures struct {
	Fixtures  []FixtureI
	Addresses []int
}

func (f *Fixtures) AddFixture(fixture FixtureI, address int) {
	f.Fixtures = append(f.Fixtures, fixture)
	f.Addresses = append(f.Addresses, address)
	// TODO: check for overlap
}

func (f *Fixtures) AddFixtures(fixture FixtureI, address ...int) Fixtures {
	list := NewFixtureList(fixture, address...)
	f.AddFixtureList(list)
	return list
}

func (f *Fixtures) AddFixtureList(fixtures Fixtures) {
	for i := range fixtures.Fixtures {
		f.AddFixture(fixtures.Fixtures[i], fixtures.Addresses[i])
	}
}

func (f Fixtures) SetValue(name string, value byte) {
	for _, fixture := range f.Fixtures {
		fixture.SetValue(name, value)
	}
}

func (f Fixtures) SetAll(value byte) {
	for _, fixture := range f.Fixtures {
		fixture.SetAll(value)
	}
}

func (f Fixtures) SetOnUpdate(onUpdate func(FixtureI), debounceTime time.Duration) {
	for _, fixture := range f.Fixtures {
		fixture.SetOnUpdate(onUpdate, debounceTime)
	}
}

func (f Fixtures) GetLastAddress() int {
	lastIndex := len(f.Addresses) - 1
	lastStartAddress := f.Addresses[lastIndex]
	lastSize := len(f.Fixtures[lastIndex].GetValues())
	return lastStartAddress + lastSize
}

func (f Fixtures) GetValues() []byte {
	res := make([]byte, f.GetLastAddress())
	for i, fixture := range f.Fixtures {
		values := fixture.GetValues()
		for j, value := range values {
			res[f.Addresses[i]+j] = value
		}
	}
	return res
}

func (f *Fixtures) Render(connection dmx.DMX) error {
	// now := time.Now()
	for i, fixture := range f.Fixtures {
		for j, value := range fixture.GetValues() {
			channel := f.Addresses[i] + j
			connection.SetChannel(channel, value)
			// log.Printf("set channel %d to value %d", channel, value)
		}
	}
	// log.Printf("rendering took %s", time.Since(now))
	return connection.Render()
}

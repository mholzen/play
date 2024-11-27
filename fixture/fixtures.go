package fixture

import (
	"log"

	"github.com/akualab/dmx"
	"github.com/mholzen/play-go/controls"
)

func NewFixtures() Fixtures {
	return make([]InstalledFixture, 0)
}

func NewFixturesFromList(constructor FixtureConstructor, address ...int) Fixtures {
	res := NewFixtures()
	for _, a := range address {
		fixture := constructor()
		res.AddFixture(fixture, a)
	}
	return res
}

func NewFixturesFromFixtures(fixtures Fixtures) Fixtures {
	res := NewFixtures()
	for _, f := range fixtures {
		res.AddFixture(f.Fixture, f.Address)
	}
	return res
}

type InstalledFixture struct {
	Fixture FixtureI
	Address int
}

func (f InstalledFixture) GetValues() []byte {
	return f.Fixture.GetValues()
}

func (f InstalledFixture) SetChannelValue(name string, value byte) {
	f.Fixture.SetChannelValue(name, value)
}

func (f InstalledFixture) SetAll(value byte) {
	f.Fixture.SetAll(value)
}

func (f InstalledFixture) SetValues(values controls.ValueMap) {
	for channel, value := range values {
		f.SetChannelValue(channel, value)
	}
}

type Fixtures []InstalledFixture

func (f *Fixtures) AddFixture(fixture FixtureI, address int) {
	*f = append(*f, InstalledFixture{fixture, address})
	// TODO: check for overlap
}

func (f *Fixtures) AddFixtures(constructor FixtureConstructor, address ...int) Fixtures {
	list := NewFixturesFromList(constructor, address...)
	f.AddFixtureList(list)
	return list
}

func (f *Fixtures) AddFixtureList(fixtures Fixtures) {
	for i := range fixtures {
		f.AddFixture(fixtures[i].Fixture, fixtures[i].Address)
	}
}

func (f Fixtures) SetChannelValue(name string, value byte) {
	for _, fixture := range f {
		fixture.Fixture.SetChannelValue(name, value)
	}
}

func (f Fixtures) SetAll(value byte) {
	for _, fixture := range f {
		fixture.Fixture.SetAll(value)
	}
}

func (f Fixtures) GetValues() []byte {
	res := make([]byte, 0)
	for _, fixture := range f {
		values := fixture.Fixture.GetValues()
		if cap(res) < fixture.Address+len(values) {
			newRes := make([]byte, fixture.Address+len(values))
			copy(newRes, res)
			res = newRes
		}

		for j, value := range values {
			res[fixture.Address+j] = value
		}
	}
	return res
}

func (f Fixtures) GetValue() FixtureValues {
	return FixtureValues{}
}

func (f Fixtures) SetValue(fixtureValues FixtureValues) {
	for fixture, values := range fixtureValues {
		fixture.SetValues(values)
	}
}

func (f *Fixtures) Render(connection dmx.DMX) error {
	for _, fixture := range *f {
		for j, value := range fixture.Fixture.GetValues() {
			channel := fixture.Address + j
			connection.SetChannel(channel, value)
			// log.Printf("set channel %d to value %d", channel, value)
		}
	}
	err := connection.Render()
	if err != nil {
		log.Printf("ERROR rendering error: %s", err)
	}
	return nil
}

func (f *Fixtures) Modulo(div, mod int) Fixtures {
	res := NewFixtures()
	for i, fixture := range *f {
		if i%div == mod {
			res = append(res, fixture)
		}
	}
	return res
}

func (f *Fixtures) Odd() Fixtures {
	return f.Modulo(2, 1)
}

func (f *Fixtures) Even() Fixtures {
	return f.Modulo(2, 0)
}

func (f Fixtures) Channel() <-chan FixtureValues {
	ch := make(chan FixtureValues)
	return ch
}

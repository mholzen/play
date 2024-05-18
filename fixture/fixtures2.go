package fixture

import (
	"time"

	"github.com/akualab/dmx"
)

func NewFixtures2() Fixtures2 {
	return make([]InstalledFixture, 0)
}

func NewFixtureList2(constructor FixtureConstructor, address ...int) Fixtures2 {
	res := NewFixtures2()
	for _, a := range address {
		fixture := constructor()
		res.AddFixture(fixture, a)
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

func (f InstalledFixture) SetValue(name string, value byte) {
	f.Fixture.SetValue(name, value)
}

func (f InstalledFixture) SetAll(value byte) {
	f.Fixture.SetAll(value)
}

func (f InstalledFixture) SetOnUpdate(onUpdate func(FixtureI), debounceTime time.Duration) {
	f.Fixture.SetOnUpdate(onUpdate, debounceTime)
}

type Fixtures2 []InstalledFixture

func (f *Fixtures2) AddFixture(fixture FixtureI, address int) {
	*f = append(*f, InstalledFixture{fixture, address})
	// TODO: check for overlap
}

func (f *Fixtures2) AddFixtures(constructor FixtureConstructor, address ...int) Fixtures2 {
	list := NewFixtureList2(constructor, address...)
	f.AddFixtureList(list)
	return list
}

func (f *Fixtures2) AddFixtureList(fixtures Fixtures2) {
	for i := range fixtures {
		f.AddFixture(fixtures[i].Fixture, fixtures[i].Address)
	}
}

func (f Fixtures2) SetValue(name string, value byte) {
	for _, fixture := range f {
		fixture.Fixture.SetValue(name, value)
	}
}

func (f Fixtures2) SetAll(value byte) {
	for _, fixture := range f {
		fixture.Fixture.SetAll(value)
	}
}

func (f Fixtures2) SetOnUpdate(onUpdate func(FixtureI), debounceTime time.Duration) {
	for _, fixture := range f {
		fixture.Fixture.SetOnUpdate(onUpdate, debounceTime)
	}
}

func (f Fixtures2) GetValues() []byte {
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

func (f *Fixtures2) Render(connection dmx.DMX) error {
	// now := time.Now()
	for _, fixture := range *f {
		for j, value := range fixture.Fixture.GetValues() {
			channel := fixture.Address + j
			connection.SetChannel(channel, value)
			// log.Printf("set channel %d to value %d", channel, value)
		}
	}
	// log.Printf("rendering took %s", time.Since(now))
	return connection.Render()
}

func (f *Fixtures2) Modulo(div, mod int) Fixtures2 {
	res := NewFixtures2()
	for i, fixture := range *f {
		if i%div == mod {
			res = append(res, fixture)
		}
	}
	return res
}

func (f *Fixtures2) Odd() Fixtures2 {
	return f.Modulo(2, 1)
}

func (f *Fixtures2) Even() Fixtures2 {
	return f.Modulo(2, 0)
}

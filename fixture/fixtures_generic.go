package fixture

import (
	"sort"

	"github.com/mholzen/play-go/controls"
)

type FixturesInterface[T FixtureI] interface {
	// AddFixture(fixture T, address int)
	// AddFixtures(constructor func() T, address ...int)
	// AddFixtureList(fixtures FixturesInterface[T])
	GetFixtures() map[int]T
	GetChannels() []string
	GetAddresses() []int
	SetValue(fixtureValues FixtureValues)
}

type FixturesGeneric[T FixtureI] map[int]*T

func NewFixturesGeneric[T FixtureI]() *FixturesGeneric[T] {
	f := make(FixturesGeneric[T])
	return &f
}

func (f *FixturesGeneric[T]) AddFixture(fixture T, address int) {
	(*f)[address] = &fixture
}

func (f *FixturesGeneric[T]) AddFixtures(constructor func() T, address ...int) {
	res := make(FixturesGeneric[T])
	for _, a := range address {
		fixture := constructor()
		res.AddFixture(fixture, a)
	}
}

// func (f *FixturesGeneric[T]) AddFixtureList(fixtures FixturesInterface[T]) {
// 	for addr, fixture := range fixtures.GetFixtures() {
// 		f.AddFixture(fixture, addr)
// 	}
// }

// func (f *FixturesGeneric[T]) Render(connection dmx.DMX) error {
// 	for _, fixture := range *f {
// 		if renderer, ok := any(fixture).(*InstalledFixture); ok {
// 			for j, value := range renderer.Fixture.GetValues() {
// 				channel := renderer.Address + j
// 				connection.SetChannel(channel, value)
// 			}
// 		}
// 	}
// 	return connection.Render()
// }

func (f *FixturesGeneric[T]) Modulo(div, mod int) FixturesGeneric[T] {
	res := make(FixturesGeneric[T])
	i := 0
	for addr, fixture := range *f {
		if i%div == mod {
			res[addr] = fixture
		}
		i++
	}
	return res
}

func (f *FixturesGeneric[T]) Odd() FixturesGeneric[T] {
	return f.Modulo(2, 1)
}

func (f *FixturesGeneric[T]) Even() FixturesGeneric[T] {
	return f.Modulo(2, 0)
}

func (f FixturesGeneric[T]) SetChannelValue(name string, value byte) {
	for _, fixture := range f {
		if setter, ok := any(fixture).(*InstalledFixture); ok {
			setter.Fixture.SetChannelValue(name, value)
		}
	}
}

func (f FixturesGeneric[T]) SetAll(value byte) {
	for _, fixture := range f {
		if setter, ok := any(fixture).(*InstalledFixture); ok {
			setter.Fixture.SetAll(value)
		}
	}
}

func (f FixturesGeneric[T]) SetValueMap(values controls.ValueMap) {
	for _, fixture := range f {
		for k, v := range values {
			(*fixture).SetChannelValue(k, v)
		}
	}
}

func (f FixturesGeneric[T]) SetFixtureValueMap(address int, values controls.ValueMap) {
	if fixture, ok := any(f[address]).(*InstalledFixture); ok {
		fixture.Fixture.SetValueMap(values)
	}
}

func (f FixturesGeneric[T]) GetValues() []byte {
	res := make([]byte, 0)
	for _, fixture := range f {
		if getter, ok := any(fixture).(*InstalledFixture); ok {
			values := getter.Fixture.GetValues()
			if cap(res) < getter.Address+len(values) {
				newRes := make([]byte, getter.Address+len(values))
				copy(newRes, res)
				res = newRes
			}
			for j, value := range values {
				res[getter.Address+j] = value
			}
		}
	}
	return res
}

func (f FixturesGeneric[T]) GetValue() FixtureValues {
	res := make(FixtureValues)
	for addr, fixture := range f {
		res[addr] = (*fixture).GetValueMap()
	}
	return res
}

func (f FixturesGeneric[T]) SetValue(fixtureValues FixtureValues) {
	for address, fixture := range f {
		(*fixture).SetValueMap(fixtureValues[address])
	}
}

func (f FixturesGeneric[T]) GetFixtureList() []*InstalledFixture {
	res := make([]*InstalledFixture, 0)
	for _, fixture := range f {
		if installed, ok := any(fixture).(*InstalledFixture); ok {
			res = append(res, installed)
		}
	}
	return res
}

func (f FixturesGeneric[T]) GetChannels() []string {
	channels := make(map[string]struct{})
	for _, fixture := range f {
		for _, channel := range (*fixture).GetChannels() {
			channels[channel] = struct{}{}
		}
	}
	keys := make([]string, 0, len(channels))
	for k := range channels {
		keys = append(keys, k)
	}
	return keys
}

func (f FixturesGeneric[T]) GetAddresses() []int {
	addresses := make([]int, 0)
	for addr := range f {
		addresses = append(addresses, addr)
	}
	return addresses
}

func (f FixturesGeneric[T]) GetFixtures() map[int]T {
	res := make(map[int]T)
	for addr, fixture := range f {
		res[addr] = *fixture
	}
	return res
}

func (f FixturesGeneric[T]) GetFixturesList() []T {
	// Get addresses and sort them
	addresses := f.GetAddresses()
	sort.Ints(addresses)

	// Build sorted list of fixtures
	res := make([]T, len(addresses))
	for i, addr := range addresses {
		res[i] = *f[addr]
	}
	return res
}

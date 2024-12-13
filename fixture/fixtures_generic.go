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
	GetValue() FixtureValues
	GetByteArray() []byte
	// Clone() FixturesInterface[T]
}

type FixturesGeneric[T FixtureI] map[int]T

type Fixtures = FixturesGeneric[Fixture]

func NewFixturesGeneric[T FixtureI]() *FixturesGeneric[T] {
	f := make(FixturesGeneric[T])
	return &f
}

func (f FixturesGeneric[T]) Clone() FixturesInterface[FixtureI] {
	res := make(FixturesGeneric[FixtureI])
	for addr, fixture := range f {
		var fixtureI FixtureI = fixture
		if fix, ok := fixtureI.(Fixture); ok {
			res[addr] = fix.Clone()
		} else {
			res[addr] = fixture
		}
	}
	return res
}

func (f *FixturesGeneric[T]) AddFixture(fixture T, address int) {
	(*f)[address] = fixture
}

func (f *FixturesGeneric[T]) AddFixtures(constructor func() T, addresses ...int) FixturesGeneric[T] {
	res := make(FixturesGeneric[T])
	for _, address := range addresses {
		fixture := constructor()
		f.AddFixture(fixture, address)
		res.AddFixture(fixture, address)
	}
	return res
}

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
	for address := range f {
		f[address].SetChannelValue(name, value)
	}
}

func (f FixturesGeneric[T]) SetAll(value byte) {
	for address := range f {
		f[address].SetAll(value)
	}
}

func (f FixturesGeneric[T]) SetValueMap(values controls.ValueMap) {
	for address := range f {
		for k, v := range values {
			f[address].SetChannelValue(k, v)
		}
	}
}

func (f FixturesGeneric[T]) SetFixtureValueMap(address int, values controls.ValueMap) {
	f[address].SetValueMap(values)
}

func (f FixturesGeneric[T]) GetValues() []byte {
	res := make([]byte, 0)
	for address, fixture := range f {
		values := fixture.GetValues()
		if cap(res) < address+len(values) {
			newRes := make([]byte, address+len(values))
			copy(newRes, res)
			res = newRes
		}
		for addr, value := range values {
			res[address+addr] = value
		}
	}
	return res
}

func (f FixturesGeneric[T]) GetValue() FixtureValues {
	res := make(FixtureValues)
	for addr, fixture := range f {
		res[addr] = fixture.GetValueMap()
	}
	return res
}

func (f FixturesGeneric[T]) SetValue(fixtureValues FixtureValues) {
	for address := range f {
		f[address].SetValueMap(fixtureValues[address])
	}
}

func (f FixturesGeneric[T]) GetChannels() []string {
	channels := make(map[string]struct{})
	for _, fixture := range f {
		for _, channel := range fixture.GetChannels() {
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

func (f FixturesGeneric[T]) GetFixtures() map[int]FixtureI {
	res := make(map[int]FixtureI)
	for addr, fixture := range f {
		res[addr] = fixture
	}
	return res
}

func (f FixturesGeneric[T]) GetFixturesList() []T { // TODO: rename to GetFixtureList()
	// Get addresses and sort them
	addresses := f.GetAddresses()
	sort.Ints(addresses)

	// Build sorted list of fixtures
	res := make([]T, len(addresses))
	for i, addr := range addresses {
		res[i] = f[addr]
	}
	return res
}

func (f FixturesGeneric[T]) GetByteArray() []byte {
	return f.GetValues()
}

func (f FixturesGeneric[T]) GetValueMap() controls.ValueMap {
	res := make(controls.ValueMap)
	for _, fixture := range f.GetFixturesList() {
		for channel, value := range fixture.GetValueMap() {
			res[channel] = value
		}
	}
	return res
}

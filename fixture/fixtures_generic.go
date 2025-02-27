package fixture

import (
	"sort"

	"github.com/mholzen/play-go/controls"
)

type Fixtures[T Fixture] interface {
	GetFixtures() map[int]T
	GetAddresses() []int
	SetValue(fixtureValues FixtureValues)
	SetChannelValue(channel string, value byte)
	GetValue() FixtureValues
	GetByteArray() []byte
	Clone() Fixtures[T]
}

type AddressableFixtures[T Fixture] map[int]T // TODO: not sure generic is useful here

type AddressableChannelFixtures = AddressableFixtures[ChannelFixture]

func NewAddressableFixtures[T Fixture]() *AddressableFixtures[T] {
	f := make(AddressableFixtures[T])
	return &f
}

func (f AddressableFixtures[T]) Clone() Fixtures[Fixture] {
	res := make(AddressableFixtures[Fixture])
	for addr, fixture := range f {
		var fixtureI Fixture = fixture
		if channelFixture, ok := fixtureI.(ChannelFixture); ok {
			res[addr] = channelFixture.Clone()
		} else {
			res[addr] = fixture
		}
	}
	return res
}

func (f *AddressableFixtures[T]) AddFixture(fixture T, address int) {
	(*f)[address] = fixture
}

func (f *AddressableFixtures[T]) AddFixtures(constructor func() T, addresses ...int) AddressableFixtures[T] {
	res := make(AddressableFixtures[T])
	for _, address := range addresses {
		fixture := constructor()
		f.AddFixture(fixture, address)
		res.AddFixture(fixture, address)
	}
	return res
}

func (f *AddressableFixtures[T]) Modulo(div, mod int) AddressableFixtures[T] {
	res := make(AddressableFixtures[T])
	i := 0
	for addr, fixture := range *f {
		if i%div == mod {
			res[addr] = fixture
		}
		i++
	}
	return res
}

func (f *AddressableFixtures[T]) Odd() AddressableFixtures[T] {
	return f.Modulo(2, 1)
}

func (f *AddressableFixtures[T]) Even() AddressableFixtures[T] {
	return f.Modulo(2, 0)
}

func (f AddressableFixtures[T]) GetChannelValue(name string) byte {
	panic("not implemented")
}

func (f AddressableFixtures[T]) SetChannelValue(name string, value byte) {
	for address := range f {
		f[address].SetChannelValue(name, value)
	}
}

func (f AddressableFixtures[T]) SetAll(value byte) {
	for address := range f {
		f[address].SetAll(value)
	}
}

func (f AddressableFixtures[T]) SetChannelValues(values controls.ChannelValues) {
	for address := range f {
		for k, v := range values {
			f[address].SetChannelValue(k, v)
		}
	}
}

func (f AddressableFixtures[T]) SetFixtureValueMap(address int, values controls.ChannelValues) {
	f[address].SetChannelValues(values)
}

func (f AddressableFixtures[T]) GetValues() []byte {
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

func (f AddressableFixtures[T]) GetValue() FixtureValues {
	res := make(FixtureValues)
	for addr, fixture := range f {
		res[addr] = fixture.GetChannelValues()
	}
	return res
}

func (f AddressableFixtures[T]) SetValue(fixtureValues FixtureValues) {
	for address := range f {
		f[address].SetChannelValues(fixtureValues[address])
	}
}

func (f AddressableFixtures[T]) GetChannels() []string {
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

func (f AddressableFixtures[T]) GetAddresses() []int {
	addresses := make([]int, 0)
	for addr := range f {
		addresses = append(addresses, addr)
	}
	return addresses
}

func (f AddressableFixtures[T]) GetFixtures() map[int]Fixture {
	res := make(map[int]Fixture)
	for addr, fixture := range f {
		res[addr] = fixture
	}
	return res
}

func (f AddressableFixtures[T]) GetFixtureList() []T {
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

func (f AddressableFixtures[T]) GetByteArray() []byte {
	return f.GetValues()
}

func (f AddressableFixtures[T]) GetChannelValues() controls.ChannelValues {
	res := make(controls.ChannelValues)
	for _, fixture := range f.GetFixtureList() {
		for channel, value := range fixture.GetChannelValues() {
			res[channel] = value
		}
	}
	return res
}

func SetChannelValues(f Fixtures[Fixture], values controls.ChannelValues) {
	for _, fixture := range f.GetFixtures() {
		fixture.SetChannelValues(values)
	}
}

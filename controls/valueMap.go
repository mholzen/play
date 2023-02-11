package controls

import "play-go/fixture"

type ValueMap map[string]byte

func (values ValueMap) ApplyTo(f fixture.FixtureI) {
	for k, v := range values {
		f.SetValue(k, v)
	}
}

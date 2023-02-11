package fixture

type FreedomPar struct {
	ModelChannels
}

func NewFreedomPar() FixtureI {
	model := NewModelChannels(
		"Freedom Par",
		[]string{
			"dimmer",
			"r", "g", "b", "a", "w", "uv",
			"strobe",
		},
	)
	return Fixture{
		Model:  model,
		Values: make([]byte, len(model.Channels)),
	}
}

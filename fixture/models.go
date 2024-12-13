package fixture

type FreedomPar struct {
	ModelChannels
}

type FixtureConstructor func() FixtureI

func NewFreedomPar() Fixture {
	model := NewModelChannels(
		"Freedom Par",
		[]string{
			"dimmer",
			"r", "g", "b", "a", "w", "uv",
			"strobe",
		},
	)
	return Fixture{
		Model:  &model,
		Values: make([]byte, len(model.Channels)),
	}
}

func NewTomeshine() Fixture {
	model := NewModelChannels(
		"Tomshine",
		[]string{
			"pan", "tilt", "speed", "dimmer", "strobe",
			"r", "g", "b", "w", "a", "uv",
		},
	)
	return Fixture{
		Model:  &model,
		Values: make([]byte, len(model.Channels)),
	}
}

func NewColorstripMini() Fixture {
	model := NewModelChannels(
		"Colorstrip Mini",
		[]string{
			"mode",
			"r", "g", "b",
		},
	)
	return Fixture{
		Model:  &model,
		Values: make([]byte, len(model.Channels)),
	}
}

// Colorstrip Mini
// DMX Channel Assignments and Values for a 4-CH DMX controller
//
// Channel 1: Static Colors
// Value 000-009: No function
// Value 010-019: Red 0-100%
// Value 020-029: Green 0-100%
// Value 030-039: Blue 0-100%
// Value 040-049: Yellow 0-100%
// Value 050-059: Magenta 0-100%
// Value 060-069: Cyan 0-100%
// Value 070-079: White 0-100%
// Value 080-089: Color Chase 1
// Value 090-099: Color Chase 2
// Value 100-109: Color Chase 3
// Value 110-119: Color Chase 4
// Value 120-129: Color Chase 5
// Value 130-139: Color Chase 6
// Value 140-149: Color Chase 7
// Value 150-159: Color Chase 8
// Value 160-169: Color Chase 9
// Value 170-179: Color Chase 10
// Value 180-189: Color Chase 11
// Value 190-199: Color Chase 12
// Value 200-209: Color Chase 13
// Value 210-219: RGB Color Mixing (Channels 2-4)
// Value 220-229: Color Fade

// Channel 2: Sound-Active
// Value 230-255: Sound-Active
// Value 000-127: Run Speed (when Ch. 1 is 080-209) Slow to fast
// Value 128-255: Sound-Active
//
// Channel 2: RGB Color Mixing
// Value 000-255: Red (when Ch. 1 is 210-219) 0-100%
// Value 000-255: Fade Speed (when Ch. 1 is 220-229) Slow to fast
//
// Channel 3: RGB Color Mixing
// Value 000-249: Strobe (when Ch. 1 is 010-119) Slow to fast
// Value 250-255: Sound-Active
// Value 000-255: Green (when Ch. 1 is 210-219) 0-100%
//
// Channel 4: RGB Color Mixing
// Value 000-255: Blue (when Ch. 1 is 210-119) 0-100%
//

func NewParCan() Fixture {
	model := NewModelChannels(
		"Battery Par Can",
		[]string{
			"dimmer",
			"r", "g", "b", "w",
			"strobe",
			"mode",
			"colorSelection",
		},
	)
	return Fixture{
		Model:  &model,
		Values: make([]byte, len(model.Channels)),
	}
}

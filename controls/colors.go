package controls

type Color struct {
	Red   byte `json:"r"`
	Green byte `json:"g"`
	Blue  byte `json:"b"`
	White byte `json:"w"`
	Amber byte `json:"a"`
	UV    byte `json:"uv"`
}

func (c Color) Values() ChannelValues {
	return ChannelValues{
		"r":  c.Red,
		"g":  c.Green,
		"b":  c.Blue,
		"w":  c.White,
		"a":  c.Amber,
		"uv": c.UV,
	}
}

type Colors map[string]Color

var AllColors Colors

func LoadColors() error {
	AllColors = make(Colors)
	for name, color := range ColorsByName {
		AllColors[name] = Color{
			Red:   color["r"],
			Green: color["g"],
			Blue:  color["b"],
			White: color["w"],
			Amber: color["a"],
			UV:    color["uv"],
		}
	}
	return nil
}

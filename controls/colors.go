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

// ObjectSum adds two colors together
func ObjectSum(a, b Color) Color {
	return Color{
		Red:   a.Red + b.Red,
		Green: a.Green + b.Green,
		Blue:  a.Blue + b.Blue,
		White: a.White + b.White,
		Amber: a.Amber + b.Amber,
		UV:    a.UV + b.UV,
	}
}

// ObjectsSum takes a slice of colors and returns their average
func ObjectsSum(colors []Color) Color {
	if len(colors) == 0 {
		return Color{}
	}

	sum := Color{}
	for _, c := range colors {
		sum = ObjectSum(sum, c)
	}

	// Calculate average
	count := byte(len(colors))
	return Color{
		Red:   sum.Red / count,
		Green: sum.Green / count,
		Blue:  sum.Blue / count,
		White: sum.White / count,
		Amber: sum.Amber / count,
		UV:    sum.UV / count,
	}
}

// Sum takes multiple colors and returns their average
func Sum(colors ...Color) Color {
	return ObjectsSum(colors)
}

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

package controls

func NewDial() *Dial {
	return &Dial{
		Value:   0,
		Channel: make(chan byte),
	}
}

type Dial struct {
	Value   byte
	Channel chan byte
}

func (d *Dial) SetValue(value byte) {
	d.Value = value
	d.Emit()
}

func (d *Dial) SetMax() {
	d.SetValue(255)
}

func (d *Dial) SetMin() {
	d.SetValue(0)
}

func (d *Dial) Toggle() {
	if d.Value <= 127 {
		d.SetMax()
	} else {
		d.SetMin()
	}
}

func (d Dial) Opposite() byte {
	x := int(d.Value) - 255
	if x < 0 {
		return byte(-x)
	} else {
		return byte(x)
	}
}

func (d Dial) Emit() {
	d.Channel <- d.Value
}

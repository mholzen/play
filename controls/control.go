package controls

type Control interface {
	Item
	GetValue() string
	SetValue() string
}

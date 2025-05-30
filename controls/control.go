package controls

type Control interface {
	Item
	GetValueString() string
	SetValueString(string) // TODO: add error handling
}

package controls

type Control interface { // TODO: use a SetStringer interface similar to fmt.Stringer
	Item
	GetValueString() string
	SetValueString(string) error
}

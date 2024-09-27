package controls

type Container interface {
	Item
	GetItem(string) (Item, error)
	Items() map[string]Item
}

func ContainerKeys(c Container) []string {
	keys := make([]string, 0, len(c.Items()))
	for k := range c.Items() {
		keys = append(keys, k)
	}
	return keys
}

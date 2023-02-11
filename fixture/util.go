package fixture

func ArrayToIndex[K comparable](array []K) map[K]int {
	r := make(map[K]int)
	for i, v := range array {
		// TODO: fail on duplicate
		r[v] = i
	}
	return r
}

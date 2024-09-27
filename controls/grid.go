package controls

import (
	"strconv"
	"strings"
)

type Grid struct {
	items [][]Item
}

func NewGrid(rows, cols int) *Grid {
	items := make([][]Item, rows)
	for i := range items {
		items[i] = make([]Item, cols)
	}
	return &Grid{items: items}
}

func (g *Grid) SetItem(row, col int, item Item) {
	g.items[row][col] = item
}

func (g *Grid) GetItem(name string) Item {
	parts := strings.Split(name, "/")
	if len(parts) != 2 {
		return nil
	}
	row, _ := strconv.Atoi(parts[0])
	col, _ := strconv.Atoi(parts[1])
	return g.items[row][col]
}

func (g *Grid) GetString() string {
	return "<multiple items>"
}

package controls

import (
	"strconv"

	"github.com/nsf/termbox-go"
)

type TermTrigger struct {
	keys             map[termbox.Key]func()
	floatKeys        map[termbox.Key]func(float64)
	input            []rune
	enteringFloatKey termbox.Key
}

func NewTermTrigger(keys map[termbox.Key]func(), floatKeys map[termbox.Key]func(float64)) *TermTrigger {
	return &TermTrigger{
		keys:      keys,
		floatKeys: floatKeys,
	}
}

func (tt *TermTrigger) Start() {
	go func() {
		err := termbox.Init()
		if err != nil {
			panic(err)
		}
		defer termbox.Close()

		for {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				val := ev.Key
				if ev.Key == 0 {
					val = termbox.Key(ev.Ch)
				}
				if callback, ok := tt.keys[val]; ok {
					callback()
				} else if _, ok := tt.floatKeys[val]; ok {
					termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
					printString(1, 1, "Enter float: ", termbox.ColorWhite, termbox.ColorBlack)
					tt.enteringFloatKey = val
					tt.input = []rune{}
				} else if tt.enteringFloatKey != 0 && ev.Ch != 0 {
					tt.input = append(tt.input, ev.Ch)
					printString(1, 10, string(tt.input), termbox.ColorWhite, termbox.ColorBlack)
				} else if tt.enteringFloatKey != 0 && ev.Key == termbox.KeyEnter {
					f, err := strconv.ParseFloat(string(tt.input), 64)
					if err == nil {
						tt.floatKeys[tt.enteringFloatKey](f)
					}
					tt.enteringFloatKey = 0
				}
			case termbox.EventError:
				panic(ev.Err)
			}
			termbox.Flush()
		}
	}()
}

func printString(x, y int, str string, fg, bg termbox.Attribute) {
	for _, c := range str {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

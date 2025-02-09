package patterns

import (
	"log"
	"testing"
	"time"

	"github.com/fogleman/ease"
	"github.com/mholzen/play-go/controls"
	"github.com/stretchr/testify/assert"
)

func Test_Transition(t *testing.T) {
	start := controls.ChannelValues{"x": 0}
	end := controls.ChannelValues{"x": 255}

	duration := time.Millisecond
	period := time.Microsecond * 100 // 10 steps over 1ms

	// Create and execute the transition
	last := start
	transition := TransitionValues(start, end, duration, ease.Linear, period, func(values controls.ChannelValues) {
		log.Printf("values: %s", values)
		last = values
	})
	transition()

	assert.Greater(t, last, start)
	assert.Equal(t, end, last)
}

package controls

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_TickPeriod(t *testing.T) {
	c := Clock{
		Bpm:          120.0,
		TicksPerBeat: 4,
		BeatsPerBar:  4,
	}
	assert.Equal(t, 125*time.Millisecond, c.TickPeriod())
}

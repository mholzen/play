package patterns

import (
	"testing"

	"github.com/mholzen/play-go/controls"
)

func Test_RainbowControls(t *testing.T) {
	rainbow := NewRainbowControls(controls.NewClock(120))
	var _ controls.Container = rainbow

}

package patterns

import (
	"testing"

	"github.com/mholzen/play/controls"
)

func TestSequenceEmitterDial(t *testing.T) {
	sequence := controls.NewSequence([]int{1, 2, 3})
	clock := controls.NewClock(120)
	sequenceEmitter := NewSequenceEmitterDial(sequence, clock)

	var _ controls.Container = sequenceEmitter
}

package patterns

import (
	"sync"
	"testing"
	"time"

	"github.com/fogleman/ease"
	"github.com/mholzen/play/controls"
	"github.com/stretchr/testify/assert"
)

func TestSequenceEmitter(t *testing.T) {
	t.Run("emits sequence values at correct intervals", func(t *testing.T) {
		values := []int{10, 20, 30}
		sequence := controls.NewSequence(values)
		clock := controls.NewClock(120) // 120 BPM

		emitter := NewSequenceEmitter(sequence, clock, controls.TriggerOnTick(2)) // Every 2 ticks

		receivedValues := []int{}
		observer := make(chan int)
		emitter.AddObserver(observer)

		go func() {
			for value := range observer {
				receivedValues = append(receivedValues, value)
			}
		}()

		clock.Start()
		time.Sleep(100 * time.Millisecond) // Let it run for a bit
		clock.Stop()

		assert.NotEmpty(t, receivedValues, "should receive values")
		if len(receivedValues) > 0 {
			assert.Equal(t, 10, receivedValues[0], "first value should be 10")
		}
	})

	t.Run("resets properly", func(t *testing.T) {
		values := []int{1, 2, 3}
		sequence := controls.NewSequence(values)
		clock := controls.NewClock(120)

		emitter := NewSequenceEmitter(sequence, clock, controls.TriggerOnTick(1))

		sequence.Inc() // Move to second value
		emitter.Reset()

		assert.Equal(t, 1, sequence.Value(), "sequence should reset to first value")
	})
}

func TestChangeEmitter(t *testing.T) {
	t.Run("only emits when value changes", func(t *testing.T) {
		source := controls.NewObservable[int]()
		emitter := NewChangeEmitter(source)

		receivedValues := []int{}
		observer := make(chan int)
		emitter.AddObserver(observer)

		go func() {
			for value := range observer {
				receivedValues = append(receivedValues, value)
			}
		}()

		// Send same value multiple times
		source.Notify(10)
		source.Notify(10)
		source.Notify(10)
		source.Notify(20) // Different value
		source.Notify(20) // Same as previous
		source.Notify(30) // Different again

		time.Sleep(10 * time.Millisecond) // Give goroutines time to process

		expected := []int{10, 20, 30}
		assert.Equal(t, expected, receivedValues, "should only emit changed values")

		emitter.Close()
	})

	t.Run("resets properly", func(t *testing.T) {
		source := controls.NewObservable[int]()
		emitter := NewChangeEmitter(source)

		source.Notify(42)
		emitter.Reset()

		assert.False(t, emitter.hasValue, "hasValue should be false after reset")
		assert.Nil(t, emitter.lastValue, "lastValue should be nil after reset")

		emitter.Close()
	})
}

func TestIntegration(t *testing.T) {
	t.Run("discretizer -> sequence -> emitter -> change emitter", func(t *testing.T) {
		discretizer := NewDiscretizer(0, 10, ease.Linear, 5)
		values := discretizer.GetValues()

		sequence := controls.NewSequence(values)

		clock := controls.NewClock(240)
		emitter := NewSequenceEmitter(sequence, clock, controls.TriggerOnTick(1))

		changeEmitter := NewChangeEmitter(emitter)

		finalValues := []int{}
		observer := make(chan int)
		changeEmitter.AddObserver(observer)

		var stopOnce sync.Once
		stopClock := func() {
			stopOnce.Do(func() {
				clock.Stop()
			})
		}

		go func() {
			for value := range observer {
				finalValues = append(finalValues, value)
				if len(finalValues) >= 10 {
					stopClock()
				}
			}
		}()

		clock.Start()
		time.Sleep(200 * time.Millisecond)
		stopClock()

		assert.NotEmpty(t, finalValues, "should receive values")
		if len(finalValues) > 0 {
			assert.Equal(t, 0, finalValues[0], "first value should be 0")
		}

		for i := 1; i < len(finalValues); i++ {
			assert.NotEqual(t, finalValues[i-1], finalValues[i],
				"change emitter should not emit duplicate consecutive values")
		}

		changeEmitter.Close()
	})
}

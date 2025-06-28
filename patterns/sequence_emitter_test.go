package patterns

import (
	"sync"
	"testing"
	"time"

	"github.com/fogleman/ease"
	"github.com/mholzen/play-go/controls"
)

func TestSequenceEmitter(t *testing.T) {
	t.Run("emits sequence values at correct intervals", func(t *testing.T) {
		values := []int{10, 20, 30}
		sequence := controls.NewSequenceT(values)
		clock := controls.NewClock(120) // 120 BPM

		emitter := NewSequenceEmitter(sequence, clock, 2) // Every 2 ticks

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

		if len(receivedValues) == 0 {
			t.Error("Expected to receive values but got none")
		}

		if receivedValues[0] != 10 {
			t.Errorf("Expected first value to be 10, got %d", receivedValues[0])
		}
	})

	t.Run("resets properly", func(t *testing.T) {
		values := []int{1, 2, 3}
		sequence := controls.NewSequenceT(values)
		clock := controls.NewClock(120)

		emitter := NewSequenceEmitter(sequence, clock, 1)

		sequence.Inc() // Move to second value
		emitter.Reset()

		if sequence.Values() != 1 {
			t.Errorf("Expected sequence to reset to first value (1), got %d", sequence.Values())
		}
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
		if len(receivedValues) != len(expected) {
			t.Fatalf("Expected %d values, got %d: %v", len(expected), len(receivedValues), receivedValues)
		}

		for i, expectedValue := range expected {
			if receivedValues[i] != expectedValue {
				t.Errorf("At index %d: expected %d, got %d", i, expectedValue, receivedValues[i])
			}
		}

		emitter.Close()
	})

	t.Run("resets properly", func(t *testing.T) {
		source := controls.NewObservable[int]()
		emitter := NewChangeEmitter(source)

		source.Notify(42)
		emitter.Reset()

		if emitter.hasValue {
			t.Error("Expected hasValue to be false after reset")
		}
		if emitter.lastValue != nil {
			t.Error("Expected lastValue to be nil after reset")
		}

		emitter.Close()
	})
}

func TestIntegration(t *testing.T) {
	t.Run("discretizer -> sequence -> emitter -> change emitter", func(t *testing.T) {
		discretizer := NewDiscretizer(0, 10, ease.Linear, 5)
		values := discretizer.GetValues()

		sequence := controls.NewSequenceT(values)

		clock := controls.NewClock(240)
		emitter := NewSequenceEmitter(sequence, clock, 1)

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

		if len(finalValues) == 0 {
			t.Error("Expected to receive values but got none")
		}

		if finalValues[0] != 0 {
			t.Errorf("Expected first value to be 0, got %d", finalValues[0])
		}

		for i := 1; i < len(finalValues); i++ {
			if finalValues[i] == finalValues[i-1] {
				t.Errorf("Change emitter should not emit duplicate consecutive values, but got %d twice", finalValues[i])
			}
		}

		changeEmitter.Close()
	})
}

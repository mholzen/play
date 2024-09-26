package controls

import (
	"sort"
	"time"
)

type TimeKeeper struct {
	times []time.Time
	cap   int
}

func NewTimeKeeper(cap int) *TimeKeeper {
	return &TimeKeeper{
		times: make([]time.Time, 0, cap),
		cap:   cap,
	}
}

func (tk *TimeKeeper) AddTime(t time.Time) {
	if len(tk.times) == tk.cap {
		tk.times = tk.times[1:]
	}
	tk.times = append(tk.times, t)
}

func (tk *TimeKeeper) GetBpm() (averageBpm, medianBpm float64) {
	if len(tk.times) < 2 {
		return 0, 0
	}
	diffs := make([]float64, len(tk.times)-1)
	for i := 0; i < len(tk.times)-1; i++ {
		diffs[i] = tk.times[i+1].Sub(tk.times[i]).Seconds()
	}
	sort.Float64s(diffs)
	average := sum(diffs) / float64(len(diffs))
	median := diffs[len(diffs)/2]
	if len(diffs)%2 == 0 {
		median = (diffs[len(diffs)/2-1] + diffs[len(diffs)/2]) / 2
	}
	averageBpm = 60 / average
	medianBpm = 60 / median
	return
}

func sum(slice []float64) float64 {
	total := 0.0
	for _, value := range slice {
		total += value
	}
	return total
}

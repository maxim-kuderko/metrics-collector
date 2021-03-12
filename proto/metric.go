package proto

import "sync"

var MetricPool = sync.Pool{New: func() interface{} { return &Metric{Values: &Values{}} }}

type Metrics map[uint64]*Metric

func (m Metrics) Reset() {
	for k := range m {
		delete(m, k)
	}
}

func (am *Metric) Add(value float64) {
	am.Values.Count++
	am.Values.Sum += value
	if value < am.Values.Min {
		am.Values.Min = value
	}
	if value > am.Values.Max {
		am.Values.Max = value
	}
}

func (am *Metric) Merge(new *Metric) {
	am.Values.Count += new.Values.Count
	am.Values.Sum += new.Values.Sum
	if new.Values.Min < am.Values.Min {
		am.Values.Min = new.Values.Min
	}
	if new.Values.Max > am.Values.Max {
		am.Values.Max = new.Values.Max
	}
}

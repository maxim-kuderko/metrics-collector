package proto

import "sync"

var MetricPool = sync.Pool{New: func() interface{} { return &Metric{Values: &Values{}} }}
var MetricRequestPool = sync.Pool{New: func() interface{} { return &MetricsRequest{} }}

type Metrics struct {
	data map[uint64]*Metric
	mu   *sync.Mutex
}

func NewMetrics() *Metrics {
	return &Metrics{
		data: map[uint64]*Metric{},
		mu:   &sync.Mutex{},
	}
}

func (m *Metrics) Reset() {
	/*m.mu.Lock()
	defer m.mu.Unlock()
	for k := range m.data {
		delete(m.data, k)
	}*/
}

func (m *Metrics) Add(new *Metric) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.data[new.Hash]
	if !ok {
		tmp := &(*new)
		tmp.Values = &(*new.Values)
		m.data[new.Hash] = tmp
		return
	}
	v.Merge(new)
}
func (m *Metrics) Data() map[uint64]*Metric {
	m.mu.Lock()
	defer m.mu.Unlock()
	defer func() {
		m.data = map[uint64]*Metric{}
	}()
	return m.data
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
	am.Values.Last = new.Values.First
}

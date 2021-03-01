package requests

import (
	"github.com/maxim-kuderko/metrics/entities"
	"sync"
)

var MetricPool = sync.Pool{New: func() interface{} { return &entities.AggregatedMetric{} }}

func Reset(m *entities.AggregatedMetric) {
	m.Hash = 0
	m.Name = ``
	m.Time = 0
	m.Tags = nil
	m.Values.Count = 0
	m.Values.Sum = 0
	m.Values.Min = 0
	m.Values.Max = 0
}

package service

import (
	"fmt"
	"github.com/maxim-kuderko/metrics-collector/internal/repositories"
	"github.com/maxim-kuderko/metrics-collector/proto"
	"github.com/spf13/viper"
	"runtime"
	"sync"
	"time"
)

type ServiceFunc func(r interface{}) (interface{}, error)

type Service struct {
	buffer *sync.Map
	shards uint64
	ticker *time.Ticker
	done   chan bool

	primaryRepo repositories.Repo
}

func NewService(p repositories.Repo, v *viper.Viper) *Service {
	s := &Service{
		buffer:      &sync.Map{},
		shards:      uint64(runtime.GOMAXPROCS(0)),
		done:        make(chan bool, 1),
		primaryRepo: p,
		ticker:      time.NewTicker(v.GetDuration(`FLUSH_INTERVAL`) * time.Millisecond),
	}
	go s.flusher()
	return s
}

func (r *Service) flusher() {
	for {
		select {
		case <-r.ticker.C:
			r.buffer.Range(func(key, value interface{}) bool {
				metrics := value.(*proto.Metrics)
				r.primaryRepo.Send(metrics)
				metrics.Reset()
				return true
			})
		case <-r.done:
			return
		}
	}
}

func (r *Service) Send(metric *proto.Metric) {
	r.send(metric)
}

func (r *Service) send(metric *proto.Metric) {
	shard := metric.Hash % r.shards
	v, _ := r.buffer.LoadOrStore(shard, proto.NewMetrics())
	v.(*proto.Metrics).Add(metric)
}

func (r *Service) Close() {
	r.done <- true
	r.buffer.Range(func(key, value interface{}) bool {
		value.(*proto.Metrics).Reset()
		return true
	})
}

func (r *Service) flush(i int) {
	v, _ := r.buffer.Load(i)
	metrics := v.(*proto.Metrics)
	if len(metrics.Data()) == 0 {
		return
	}
	if err := r.primaryRepo.Send(metrics); err != nil {
		fmt.Println(err)
	}

}

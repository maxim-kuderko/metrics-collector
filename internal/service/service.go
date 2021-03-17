package service

import (
	"fmt"
	"github.com/maxim-kuderko/metrics-collector/internal/repositories"
	"github.com/maxim-kuderko/metrics-collector/proto"
	"github.com/spf13/viper"
	"sync"
	"time"
)

type ServiceFunc func(r interface{}) (interface{}, error)

type Service struct {
	buffer         *sync.Map
	shards         uint64
	ticker         *time.Ticker
	done           chan bool
	wg             sync.WaitGroup
	flushSemaphore chan struct{}

	primaryRepo repositories.Repo
}

func NewService(p repositories.Repo, v *viper.Viper) *Service {
	s := &Service{
		buffer:         &sync.Map{},
		shards:         v.GetUint64(`SHARDS`),
		done:           make(chan bool, 1),
		flushSemaphore: make(chan struct{}, v.GetInt(`SHARDS`)*2),
		primaryRepo:    p,
		ticker:         time.NewTicker(v.GetDuration(`FLUSH_INTERVAL`) * time.Millisecond),
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
	r.wg.Wait()
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

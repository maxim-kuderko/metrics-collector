package service

import (
	"fmt"
	"github.com/maxim-kuderko/metrics-collector/internal/repositories"
	metricsEnt "github.com/maxim-kuderko/metrics/entities"
	"github.com/spf13/viper"
	"sync"
	"time"
)

type ServiceFunc func(r interface{}) (interface{}, error)

type Service struct {
	buffer         []metricsEnt.Metrics
	mu             []*sync.Mutex
	ticker         *time.Ticker
	done           chan bool
	wg             sync.WaitGroup
	flushSemaphore chan struct{}

	primaryRepo repositories.Repo
}

var MetricsPool = sync.Pool{New: newBuff()}

func newBuff() func() interface{} {
	return func() interface{} {
		return metricsEnt.Metrics{}
	}
}

func NewService(p repositories.Repo, v *viper.Viper) *Service {
	buff := make([]metricsEnt.Metrics, 0, v.GetInt(`SHARDS`))
	mu := make([]*sync.Mutex, 0, v.GetInt(`SHARDS`))
	for i := 0; i < v.GetInt(`SHARDS`); i++ {
		buff = append(buff, MetricsPool.Get().(metricsEnt.Metrics))
		mu = append(mu, &sync.Mutex{})
	}
	s := &Service{
		buffer:         buff,
		mu:             mu,
		done:           make(chan bool, 1),
		flushSemaphore: make(chan struct{}, v.GetInt(`SHARDS`)*2),
		primaryRepo:    p,
		ticker:         time.NewTicker(v.GetDuration(`FLUSH_INTERVAL`) * time.Millisecond),
	}
	go s.flusher()
	return s
}

func (s *Service) flusher() {
	for {
		select {
		case <-s.ticker.C:
			for i, mu := range s.mu {
				mu.Lock()
				s.flush(i)
				mu.Unlock()
			}
		case <-s.done:
			return
		}
	}
}

func (r *Service) Send(metrics []metricsEnt.AggregatedMetric) {
	for _, m := range metrics {
		r.send(m)
	}
}

func (r *Service) send(metric metricsEnt.AggregatedMetric) {
	shard := metric.Hash % uint64(len(r.mu))
	r.mu[shard].Lock()
	defer r.mu[shard].Unlock()
	v, ok := r.buffer[shard][metric.Hash]
	if !ok {
		r.buffer[shard][metric.Hash] = &metric
		v = &metric
		return
	}
	v.Merge(&metric)
}

func (r *Service) Close() {
	r.done <- true
	for i, mu := range r.mu {
		mu.Lock()
		r.flush(i)
		mu.Unlock()
	}
	r.wg.Wait()
}

func (r *Service) flush(i int) {
	if len(r.buffer[i]) == 0 {
		return
	}
	r.wg.Add(1)
	tmp := r.buffer[i]
	r.buffer[i] = MetricsPool.Get().(metricsEnt.Metrics)
	r.flushSemaphore <- struct{}{}
	go func() {
		defer func() {
			<-r.flushSemaphore
			tmp.Reset()
			MetricsPool.Put(tmp)
			r.wg.Done()
		}()
		if err := r.primaryRepo.Send(tmp); err != nil {
			fmt.Println(err)
		}
	}()
}

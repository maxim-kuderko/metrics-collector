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
	buffer         []proto.Metrics
	mu             []*sync.Mutex
	ticker         *time.Ticker
	done           chan bool
	wg             sync.WaitGroup
	flushSemaphore chan struct{}

	primaryRepo repositories.Repo
}

func NewService(p repositories.Repo, v *viper.Viper) *Service {
	buff := make([]proto.Metrics, 0, v.GetInt(`SHARDS`))
	mu := make([]*sync.Mutex, 0, v.GetInt(`SHARDS`))
	for i := 0; i < v.GetInt(`SHARDS`); i++ {
		buff = append(buff, proto.Metrics{})
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

func (r *Service) flusher() {
	for {
		select {
		case <-r.ticker.C:
			for i, mu := range r.mu {
				mu.Lock()
				r.flush(i)
				mu.Unlock()
			}
		case <-r.done:
			return
		}
	}
}

func (r *Service) Send(metric *proto.Metric) {
	r.send(metric)
}

func (r *Service) send(metric *proto.Metric) {
	shard := metric.Hash % uint64(len(r.mu))
	r.mu[shard].Lock()
	defer r.mu[shard].Unlock()
	v, ok := r.buffer[shard][metric.Hash]
	if !ok {
		tmp := &(*metric)
		tmp.Values = &(*metric.Values)
		r.buffer[shard][metric.Hash] = tmp
		return
	}
	v.Merge(metric)
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
	r.buffer[i] = proto.Metrics{}
	r.flushSemaphore <- struct{}{}
	go func() {
		defer func() {
			<-r.flushSemaphore
			r.wg.Done()
		}()
		if err := r.primaryRepo.Send(tmp); err != nil {
			fmt.Println(err)
		}
	}()
}

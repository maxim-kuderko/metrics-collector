package service

import (
	"fmt"
	"github.com/maxim-kuderko/metrics-collector/internal/repositories"
	"github.com/maxim-kuderko/metrics-collector/proto"
	"github.com/spf13/viper"
	"runtime"
	"time"
)

type ServiceFunc func(r interface{}) (interface{}, error)

type Service struct {
	buffer []*proto.Metrics
	shards uint64
	ticker *time.Ticker
	done   chan bool

	primaryRepo repositories.Repo
}

func NewService(p repositories.Repo, v *viper.Viper) *Service {
	shards := runtime.GOMAXPROCS(0)
	buff := make([]*proto.Metrics, 0, shards)
	for i := 0; i < shards; i++ {
		buff = append(buff, proto.NewMetrics())
	}
	var s = &Service{
		buffer:      buff,
		shards:      uint64(shards),
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
			for _, b := range r.buffer {
				r.primaryRepo.Send(b)
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
	r.buffer[metric.Hash%r.shards].Add(metric)
}

func (r *Service) Close() {
	r.done <- true
	for _, buff := range r.buffer {
		r.primaryRepo.Send(buff)
	}
}

func (r *Service) flush(i int) {
	if err := r.primaryRepo.Send(r.buffer[i]); err != nil {
		fmt.Println(err)
	}
}

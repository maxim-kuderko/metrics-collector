package repositories

import (
	"fmt"
	"github.com/maxim-kuderko/metrics-collector/proto"
	"go.uber.org/atomic"
	"time"
)

type Stdout struct {
	c *atomic.Int64
}

func (s Stdout) Send(r proto.Metrics) error {
	s.c.Add(int64(len(r)))
	return nil
}

func NewStdout() Repo {
	s := &Stdout{c: atomic.NewInt64(0)}
	go func() {
		w := 1
		t := time.NewTicker(time.Second * time.Duration(w))
		for range t.C {
			fmt.Println(fmt.Sprintf("%0.2fm req/sec ", float64(s.c.Swap(0))/1000000/float64(w)))
		}
	}()
	return s
}

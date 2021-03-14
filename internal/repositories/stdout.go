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
	c := int64(0)
	for _, m := range r {
		c += m.Values.Count
	}
	s.c.Add(c)
	return nil
}

func NewStdout() Repo {
	s := &Stdout{c: atomic.NewInt64(0)}
	go func() {
		w := 3
		t := time.NewTicker(time.Second * time.Duration(w))
		for range t.C {
			fmt.Println(fmt.Sprintf("%0.2fm req/sec ", float64(s.c.Swap(0))/1000000/float64(w)))
		}
	}()
	return s
}

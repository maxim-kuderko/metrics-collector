package repositories

import (
	jsoniter "github.com/json-iterator/go"
	metricsEnt "github.com/maxim-kuderko/metrics/entities"
	"os"
)

type Stdout struct {
}

func (s Stdout) Send(r metricsEnt.Metrics) error {
	for _, m := range r {
		jsoniter.ConfigFastest.NewEncoder(os.Stdout).Encode(m)
	}
	return nil
}

func NewStdout() Repo {
	return &Stdout{}
}

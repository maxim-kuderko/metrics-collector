package repositories

import (
	"github.com/maxim-kuderko/metrics-collector/proto"
)

type Repo interface {
	Send(r *proto.Metrics) error
}

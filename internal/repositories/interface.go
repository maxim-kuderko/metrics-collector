package repositories

import (
	metricsEnt "github.com/maxim-kuderko/metrics/entities"
)

type Repo interface {
	Send(r metricsEnt.Metrics) error
}

package instrumentation

import (
	"sync"

	"github.com/cabify/gotoprom"
	"github.com/prometheus/client_golang/prometheus"
)

type ServiceLabels struct {
	Action  string `label:"action"`
	Success bool   `label:"success"`
}

type Metrics struct {
	Service struct {
		Total   func(ServiceLabels) prometheus.Gauge     `name:"total" help:"Total amount of service calls"`
		Latency func(ServiceLabels) prometheus.Histogram `name:"latency" help:"Latency service calls" buckets:""`
	} `namespace:"service"`
}

var (
	once sync.Once
	m    Metrics
)

func NewMetrics() (Metrics, error) {
	var err error
	once.Do(func() {
		err = gotoprom.Init(&m, "racers_backend")
	})

	return m, err
}

package tracer

import (
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/reporter/http"
)

const zipkinEndpoint = "/api/v2/spans"

func New(serviceName, baseUrl string) (*zipkin.Tracer, error) {
	// Local endpoint represent the local service information
	localEndpoint := &model.Endpoint{ServiceName: serviceName}

	// We will record 100% (1.00) of traces.
	sampler, err := zipkin.NewCountingSampler(1)
	if err != nil {
		return nil, err
	}

	return zipkin.NewTracer(
		http.NewReporter(baseUrl+zipkinEndpoint),
		zipkin.WithSampler(sampler),
		zipkin.WithLocalEndpoint(localEndpoint),
	)
}

func NewNoop() *zipkin.Tracer {
	t, _ := zipkin.NewTracer(nil)

	return t
}

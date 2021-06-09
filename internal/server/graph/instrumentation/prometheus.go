package instrumentation

import (
	"context"
	"time"

	"github.com/99designs/gqlgen/graphql"
	prometheusclient "github.com/prometheus/client_golang/prometheus"
)

const (
	operationType = "operation_type"
	operationName = "operation_name"
	objectName    = "object_name"
	objectField   = "object_field"
)

type Prometheus struct {
	requestCounter  *prometheusclient.CounterVec
	resolverCounter *prometheusclient.CounterVec
	requestLatency  *prometheusclient.HistogramVec
	resolverLatency *prometheusclient.HistogramVec
}

var _ interface {
	graphql.HandlerExtension
	graphql.OperationInterceptor
	graphql.ResponseInterceptor
	graphql.FieldInterceptor
} = Prometheus{}

// RegisterPrometheus adds observability to gqlgen
func NewPrometheus(registerer prometheusclient.Registerer, namespace string) Prometheus {
	var prom Prometheus

	prom.requestCounter = prometheusclient.NewCounterVec(
		prometheusclient.CounterOpts{
			Name:      "graphql_request",
			Help:      "Total number of requests on the graphql server.",
			Namespace: namespace,
		},
		[]string{operationType, operationName},
	)

	prom.requestLatency = prometheusclient.NewHistogramVec(
		prometheusclient.HistogramOpts{
			Name:      "graphql_request_latency",
			Help:      "The time taken to resolve a request by graphql server.",
			Namespace: namespace,
		},
		[]string{operationType, operationName},
	)

	prom.resolverCounter = prometheusclient.NewCounterVec(
		prometheusclient.CounterOpts{
			Name:      "graphql_field",
			Help:      "Total number of requests to a field.",
			Namespace: namespace,
		},
		[]string{objectName, objectField},
	)

	prom.resolverLatency = prometheusclient.NewHistogramVec(
		prometheusclient.HistogramOpts{
			Name:      "graphql_field_latency",
			Help:      "The time taken to resolve a field.",
			Namespace: namespace,
		},
		[]string{objectName, objectField},
	)

	registerer.MustRegister(
		prom.requestCounter,
		prom.resolverCounter,
		prom.requestLatency,
		prom.resolverLatency,
	)

	return prom
}

// ExtensionName is a required method to define gqlgen interceptor
func (p Prometheus) ExtensionName() string {
	return "Prometheus"
}

// Validate is a required method to define gqlgen interceptor
func (p Prometheus) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

// InterceptOperation adds one to request counter
func (p Prometheus) InterceptOperation(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	if ol := operationLabels(ctx); ol != nil {
		p.requestCounter.With(ol).Inc()
	}

	return next(ctx)
}

// InterceptResponse registers the amount of time taken to resolve a request
func (p Prometheus) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	if ol := operationLabels(ctx); ol != nil {
		oc := graphql.GetOperationContext(ctx)
		p.requestLatency.With(ol).Observe(float64(time.Since(oc.Stats.OperationStart).Milliseconds()))
	}

	return next(ctx)
}

// InterceptField adds one to field counter and registers the amount of time taken to resolve a field if it is a resolver
func (p Prometheus) InterceptField(ctx context.Context, next graphql.Resolver) (interface{}, error) {
	fc := graphql.GetFieldContext(ctx)
	labels := prometheusclient.Labels{
		objectName:  fc.Field.ObjectDefinition.Name,
		objectField: fc.Field.Name,
	}

	p.resolverCounter.With(labels).Inc()

	observerStart := time.Now()
	res, err := next(ctx)

	// If is not a resolver we don't want to know latency
	if fc.IsResolver {
		p.resolverLatency.With(labels).
			Observe(float64(time.Since(observerStart).Nanoseconds()))
	}
	return res, err
}

func operationLabels(ctx context.Context) prometheusclient.Labels {
	opCtx := graphql.GetOperationContext(ctx)
	if opCtx.Operation == nil {
		return nil
	}
	collected := graphql.CollectFields(opCtx, opCtx.Operation.SelectionSet, []string{string(opCtx.Operation.Operation)})
	if len(collected) == 0 {
		return nil
	}

	return prometheusclient.Labels{
		operationType: string(opCtx.Operation.Operation),
		operationName: collected[0].Name,
	}
}

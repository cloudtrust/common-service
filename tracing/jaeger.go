package tracing

import (
	"context"
	"fmt"
	"io"
	"net/http"

	cs "github.com/cloudtrust/common-service"
	opentracing "github.com/opentracing/opentracing-go"
	otag "github.com/opentracing/opentracing-go/ext"
	jaeger "github.com/uber/jaeger-client-go/config"
)

// Finisher interface
type Finisher interface {
	Finish()
}

type noopFinisher struct{}

func (f *noopFinisher) Finish() {}

// OpentracingClient used for Jaeger
type OpentracingClient interface {
	TryStartSpanWithTag(ctx context.Context, operationName, tagName, tagValue string) (context.Context, Finisher)
	MakeEndpointTracingMW(operationName string) cs.Middleware
	MakeHTTPTracingMW(componentName, operationName string) func(http.Handler) http.Handler
	Close()
}

type noopOpentracingClient struct {
}

func (o *noopOpentracingClient) TryStartSpanWithTag(ctx context.Context, operationName, tagName, tagValue string) (context.Context, Finisher) {
	return ctx, &noopFinisher{}
}
func (o *noopOpentracingClient) MakeEndpointTracingMW(operationName string) cs.Middleware {
	return func(next cs.Endpoint) cs.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			return next(ctx, request)
		}
	}
}
func (o *noopOpentracingClient) MakeHTTPTracingMW(componentName, operationName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}
func (o *noopOpentracingClient) Close() {}

type internalOpentracingClient struct {
	Tracer opentracing.Tracer
	closer io.Closer
}

// TryStartSpanWithTag tries to span with a given tag
func (o *internalOpentracingClient) TryStartSpanWithTag(ctx context.Context, operationName, tagName, tagValue string) (context.Context, Finisher) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span = o.Tracer.StartSpan(operationName, opentracing.ChildOf(span.Context()))
		span.SetTag(tagName, tagValue)

		return opentracing.ContextWithSpan(ctx, span), span
	}
	return ctx, nil
}

// MakeEndpointTracingMW makes a middleware that handle the tracing with jaeger.
func (o *internalOpentracingClient) MakeEndpointTracingMW(operationName string) cs.Middleware {
	return func(next cs.Endpoint) cs.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			if span := opentracing.SpanFromContext(ctx); span != nil {
				span = o.Tracer.StartSpan(operationName, opentracing.ChildOf(span.Context()))
				defer span.Finish()

				span.SetTag("correlation_id", ctx.Value("correlation_id").(string))

				ctx = opentracing.ContextWithSpan(ctx, span)
			}
			return next(ctx, request)
		}
	}
}

// MakeHTTPTracingMW try to extract an existing span from the HTTP headers. It it exists, we
// continue the span, if not we create a new one.
func (o *internalOpentracingClient) MakeHTTPTracingMW(componentName, operationName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var sc, err = o.Tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))

			var span opentracing.Span
			if err != nil {
				span = o.Tracer.StartSpan(operationName)
			} else {
				span = o.Tracer.StartSpan(operationName, opentracing.ChildOf(sc))
			}
			defer span.Finish()

			// Set tags.
			otag.Component.Set(span, componentName)
			span.SetTag("transport", "http")
			otag.SpanKindRPCServer.Set(span)

			next.ServeHTTP(w, r.WithContext(opentracing.ContextWithSpan(r.Context(), span)))
		})
	}
}

// Close the jaeger client
func (o *internalOpentracingClient) Close() {
	o.closer.Close()
}

// CreateJaegerClient creates an opentracing Jaerger client
// For its configuration, parameter names are built with the given prefix, then a dash symbol, then one of these suffixes:
// sampler-type, sampler-param, sampler-host-port, reporter-logspan, write-interval
// If a parameter exists only named with the given prefix and if its value if false, the OpentracingClient
// will be a inactive one (Noop)
func CreateJaegerClient(v cs.Configuration, prefix string, componentName string) (OpentracingClient, error) {
	if !v.GetBool(prefix) {
		return &noopOpentracingClient{}, nil
	}
	jaegerConfig := jaeger.Configuration{
		Disabled: false,
		Sampler: &jaeger.SamplerConfig{
			Type:              v.GetString(prefix + "-sampler-type"),
			Param:             v.GetFloat64(prefix + "-sampler-param"),
			SamplingServerURL: fmt.Sprintf("http://%s", v.GetString(prefix+"-sampler-host-port")),
		},
		Reporter: &jaeger.ReporterConfig{
			LogSpans:            v.GetBool(prefix + "-reporter-logspan"),
			BufferFlushInterval: v.GetDuration(prefix + "-write-interval"),
		},
	}

	tracer, closer, err := jaegerConfig.New(componentName)
	if err != nil {
		return nil, err
	}
	return &internalOpentracingClient{
		Tracer: tracer,
		closer: closer,
	}, nil
}

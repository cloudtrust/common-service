package tracing

import (
	"fmt"
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	jaeger "github.com/uber/jaeger-client-go/config"
)

// OpentracingClient used for Jaeger and closable
type OpentracingClient struct {
	Tracer opentracing.Tracer
	closer io.Closer
}

// Close the jaeger client
func (o *OpentracingClient) Close() {
	o.closer.Close()
}

// CreateJaegerClient creates an opentracing Jaerger client
func CreateJaegerClient(v *viper.Viper, prefix string, componentName string) (*OpentracingClient, error) {
	jaegerConfig := jaeger.Configuration{
		Disabled: !v.GetBool(prefix),
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
	return &OpentracingClient{
		Tracer: tracer,
		closer: closer,
	}, nil
}

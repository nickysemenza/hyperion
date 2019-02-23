package tracing

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

//InitTracer starts a global jaeger tracer
func InitTracer(serverAddress, serviceName string) {
	if serverAddress == "" {
		log.Info("tracing is not enabled")
		opentracing.SetGlobalTracer(opentracing.NoopTracer{})
		return
	}
	// sender := transport.NewHTTPTransport(
	// 	"localhost:6831",
	// 	transport.HTTPBatchSize(1),
	// )

	sender, err := jaeger.NewUDPTransport(serverAddress, 0)
	if err != nil {
		log.Fatal(err)
	}
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1.0, // sample all traces
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, _, _ := cfg.New(serviceName,
		jaegercfg.Reporter(jaeger.NewRemoteReporter(
			sender,
			jaeger.ReporterOptions.BufferFlushInterval(1*time.Second),
			jaeger.ReporterOptions.Logger(jaegerlog.StdLogger),
		)))

	// defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
	log.Info("tracing enabled", serverAddress, serviceName)

}

//GinMiddleware is a gin middleware for initializing a trace via a HTTP request
func GinMiddleware(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(ctx, "request: "+c.Request.Method+" "+c.Request.URL.Path)
		ext.HTTPMethod.Set(span, c.Request.Method)
		ext.HTTPUrl.Set(span, c.Request.URL.Path)
		defer ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))
		defer span.Finish()
		c.Set("ctx", ctx)
		if sc, ok := span.Context().(jaeger.SpanContext); ok {
			c.Writer.Header().Set("x-trace-id", sc.TraceID().String())
		}

		c.Next()
	}
}

type stringTagName string

// Set adds a string tag to the `span`
func (tag stringTagName) Set(span opentracing.Span, value string) {
	span.SetTag(string(tag), value)
}

//SetError sets the error tag to true, and logs the error.
func SetError(span opentracing.Span, err error) {
	span.SetTag("error", true)
	span.LogKV("error", err.Error())
}

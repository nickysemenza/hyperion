package tracing

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"

	"contrib.go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
)

//InitTracer starts a global jaeger tracer
func InitTracer(serverAddress, serviceName string) {
	je, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint: serverAddress,
		// CollectorEndpoint:      collectorEndpointURI,
		ServiceName: serviceName,
	})
	if err != nil {
		log.Fatalf("Failed to create the Jaeger exporter: %v", err)
	}

	// And now finally register it as a Trace Exporter
	trace.RegisterExporter(je)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	log.Info("tracing enabled", serverAddress, serviceName)

}

//GinMiddleware is a gin middleware for initializing a trace via a HTTP request
func GinMiddleware(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := trace.StartSpan(ctx, "request: "+c.Request.Method+" "+c.Request.URL.Path)
		// ext.HTTPMethod.Set(span, c.Request.Method)
		// ext.HTTPUrl.Set(span, c.Request.URL.Path)
		// defer ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))
		defer span.End()
		c.Set("ctx", ctx)
		// if sc, ok := span.Context().(jaeger.SpanContext); ok {
		// 	c.Writer.Header().Set("x-trace-id", sc.TraceID().String())
		// }

		c.Next()
	}
}

type stringTagName string

// Set adds a string tag to the `span`
func (tag stringTagName) Set(span *trace.Span, value string) {
	span.AddAttributes(trace.StringAttribute(string(tag), value))
}

//SetError sets the error tag to true, and logs the error.
func SetError(span *trace.Span, err error) {
	span.AddAttributes(trace.BoolAttribute("error", true))
	span.Annotate([]trace.Attribute{
		trace.StringAttribute("error", err.Error()),
	}, "error")
}

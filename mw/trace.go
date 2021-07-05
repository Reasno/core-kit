package mw

import (
	"context"

	"github.com/DoNewsCode/core/contract"
	"github.com/DoNewsCode/core/key"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/opentracing"
	stdtracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// LabeledTraceServer returns a LabeledMiddleware that wraps the `next` Endpoint in an
// OpenTracing Span. The name of the operation is defined by contract.Keyer.
func LabeledTraceServer(tracer stdtracing.Tracer, keyer contract.Keyer) LabeledMiddleware {
	return func(method string, endpoint endpoint.Endpoint) endpoint.Endpoint {
		name := key.KeepOdd(keyer).Key(".", "method", method)
		return TraceServer(tracer, name)(endpoint)
	}
}

// TraceServer returns a Middleware that wraps the `next` Endpoint in an
// OpenTracing Span called `operationName`.
//
// If `ctx` already has a Span, it is re-used and the operation name is
// overwritten. If `ctx` does not yet have a Span, one is created here.
func TraceServer(tracer stdtracing.Tracer, operationName string, opts ...opentracing.EndpointOption) endpoint.Middleware {
	opts = append(opts, opentracing.WithTags(map[string]interface{}{
		ext.SpanKindRPCServer.Key: ext.SpanKindRPCServer.Value,
	}))
	return traceWithTags(tracer, operationName, opts...)
}

// TraceClient returns a Middleware that wraps the `next` Endpoint in an
// OpenTracing Span called `operationName`.
//
// If `ctx` already has a Span, it is re-used and the operation name is
// overwritten. If `ctx` does not yet have a Span, one is created here.
func TraceClient(tracer stdtracing.Tracer, operationName string, opts ...opentracing.EndpointOption) endpoint.Middleware {
	opts = append(opts, opentracing.WithTags(map[string]interface{}{
		ext.SpanKindRPCServer.Key: ext.SpanKindRPCClient.Value,
	}))
	return traceWithTags(tracer, operationName, opts...)
}

// TraceConsumer returns a Middleware that wraps the `next` Endpoint in an
// OpenTracing Span called `operationName`.
//
// If `ctx` already has a Span, it is re-used and the operation name is
// overwritten. If `ctx` does not yet have a Span, one is created here.
func TraceConsumer(tracer stdtracing.Tracer, operationName string, opts ...opentracing.EndpointOption) endpoint.Middleware {
	opts = append(opts, opentracing.WithTags(map[string]interface{}{
		ext.SpanKindRPCServer.Key: ext.SpanKindConsumer.Value,
	}))
	return traceWithTags(tracer, operationName, opts...)
}

// TraceProducer returns a Middleware that wraps the `next` Endpoint in an
// OpenTracing Span called `operationName`.
func TraceProducer(tracer stdtracing.Tracer, operationName string, opts ...opentracing.EndpointOption) endpoint.Middleware {
	opts = append(opts, opentracing.WithTags(map[string]interface{}{
		ext.SpanKindRPCServer.Key: ext.SpanKindProducer.Value,
	}))
	return traceWithTags(tracer, operationName, opts...)
}

func traceWithTags(tracer stdtracing.Tracer, operationName string, opts ...opentracing.EndpointOption) endpoint.Middleware {
	opts = append(opts, opentracing.WithTagsFunc(func(ctx context.Context) stdtracing.Tags {
		var tags = make(stdtracing.Tags)
		transport, _ := ctx.Value(contract.TransportKey).(string)
		tags["transport"] = transport
		requestUrl, _ := ctx.Value(contract.RequestUrlKey).(string)
		tags["request.url"] = requestUrl
		if tenant, ok := ctx.Value(contract.TenantKey).(contract.Tenant); ok {
			for k, v := range tenant.KV() {
				tags[k] = v
			}
		}
		return tags
	}))
	return opentracing.TraceEndpoint(tracer, operationName, opts...)
}

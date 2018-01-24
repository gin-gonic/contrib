// Package gintracing provides requests tracing functional using opentracing specification.
//
// See https://github.com/opentracing/opentracing-go for more information
package gintracing

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

const spanContextKey = "span"

// AbortOnErrors determines middleware behavior on errors.
// If set to true middleware aborts request with http.InternalServerError.
var AbortOnErrors = false

var ErrSpanNotFound = errors.New("span was not found in context")

// NewSpan returns gin.HandlerFunc (middleware) that starts a new span and injects it to request context.
//
// It calls ctx.Next() to measure execution time of all following handlers.
func NewSpan(tracer opentracing.Tracer, operationName string, opts ...opentracing.StartSpanOption) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		span := tracer.StartSpan(operationName, opts...)
		ctx.Set(spanContextKey, span)
		defer span.Finish()

		ctx.Next()
	}
}

// ParentSpanReferenceFunc determines how to reference parent span
//
// See opentracing.SpanReferenceType
type ParentSpanReferenceFunc func(opentracing.SpanContext) opentracing.StartSpanOption

// SpanFromHeaders returns gin.HandlerFunc (middleware) that extracts parent span data from HTTP headers and
// starts a new span referenced to parent with ParentSpanReferenceFunc.
//
// It calls ctx.Next() to measure execution time of all following handlers.
func SpanFromHeaders(tracer opentracing.Tracer, operationName string,
	psr ParentSpanReferenceFunc, advancedOpts ...opentracing.StartSpanOption) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		spanContext, err := tracer.Extract(opentracing.TextMap, opentracing.HTTPHeadersCarrier(ctx.Request.Header))
		if err != nil {
			if AbortOnErrors {
				ctx.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}

		opts := append([]opentracing.StartSpanOption{psr(spanContext)}, advancedOpts...)

		span := tracer.StartSpan(operationName, opts...)
		ctx.Set(spanContextKey, span)
		defer span.Finish()

		ctx.Next()
	}
}

// SpanFromContext returns gin.HandlerFunc (middleware) that extracts parent span from request context
// and starts a new span as child of parent span.
//
// It calls ctx.Next() to measure execution time of all following handlers.
func SpanFromContext(tracer opentracing.Tracer, operationName string,
	advancedOpts ...opentracing.StartSpanOption) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var opts []opentracing.StartSpanOption
		parentSpanI, _ := ctx.Get(spanContextKey)
		if parentSpan, typeOk := parentSpanI.(opentracing.Span); parentSpan != nil && typeOk {
			opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
		} else {
			if AbortOnErrors {
				ctx.AbortWithError(http.StatusInternalServerError, ErrSpanNotFound)
			}
			return
		}
		opts = append(opts, advancedOpts...)

		span := tracer.StartSpan(operationName, opts...)
		ctx.Set(spanContextKey, span)
		defer span.Finish()

		ctx.Next()
	}
}

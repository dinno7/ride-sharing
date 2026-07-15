package echocustoms

import (
	echootel "github.com/labstack/echo-opentelemetry"
	"github.com/labstack/echo/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func TracingMiddlewareWithName(spanName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			tracer, err := echo.ContextGet[trace.Tracer](c, echootel.TracerKey)
			if err != nil {
				return err
			}
			// Use the request's existing context as the parent
			ctx := c.Request().Context()

			// Name your span however you like
			route := c.Request().URL

			if spanName == "" {
				spanName := route.Path
				if spanName == "" {
					spanName = c.Path()
				}
			}

			// Start span and replace request context with it
			ctx, span := tracer.Start(
				ctx, spanName,
				trace.WithSpanKind(trace.SpanKindServer),
				trace.WithAttributes(
					attribute.String("http.method", c.Request().Method),
					attribute.String("http.route", spanName),
					attribute.String("http.target", c.Request().URL.Path),
					attribute.String("http.scheme", route.Scheme),
				),
			)
			defer span.End()

			c.SetRequest(c.Request().WithContext(ctx))

			err = next(c)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "")
			}

			return err
		}
	}
}

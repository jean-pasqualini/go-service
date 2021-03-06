package mid

import (
	"context"
	"github.com/jean-pasqualini/go-service/foundation/web"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"runtime/debug"
)

// Panics recovers from panics and convert the panic to an error so it is reported in Metrics and handled in Errors.
// MIddleware factory
func Panics(log *log.Logger) web.Middleware {
	// Decorator
	return func(handler web.Handler) web.Handler {
		// Decorated
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.mid.panics")
			defer span.End()

			// If the context is missing this value, request the service to be shutdown gracefully.
			v, ok := ctx.Value(web.KeyValues).(*web.Values)
			if !ok {
				return web.NewShutdownError("web value missing from context")
			}

			// Defer a function to recover from a panic and set the err return variable after the fact
			defer func() {
				if r := recover(); r != nil {
					err = errors.Errorf("PANIC     : %v", r)

					// Log the Go stack trace for this panic'd goroutine.
					log.Printf("%s : PANIC     :\n%s", v.TraceId, debug.Stack())
				}
			}()

			// Call the Handler and set its return value in the err variable
			return handler(ctx, w, r)
		}
	}
}

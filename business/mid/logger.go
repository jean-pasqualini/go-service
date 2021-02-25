package mid

import (
	"context"
	"github.com/jean-pasqualini/go-service/foundation/web"
	"log"
	"net/http"
	"time"
)

// Logger writes some information about the request to the logs in the format:
// TraceID : (200) GET /foo -> IP ADDR (latency)
// ---------------------------------------------------------------------------------------------------------------------
// Create a new middleware func() with a logger available in his scope
// A middleware is a handler decorator,
// It takes an handler as first parameter (generic handler) and wrap it with an specialized handler (loggerHandler here)
// The handler returned is compatible with httpTreeMux
// -- ServerHandler
// ------ http.TreeMuxHandler
// -------- web.Handler - build by web.Middleware (decorate) when it's executed  (run the inner handler in his scope)
// ------------- web.Handler - build by web.Midleware (decorate) when it's executed (run the inner handler in his scope)
// ----------------- web.Handler - your specific handler
func Logger(log *log.Logger) web.Middleware {

	// It will decorate the handler passed with a custom behavior before and/or after have called the inner handler.
	middleware := func (innerHandler web.Handler) web.Handler {

		// This is a decorated handler.
		outerHandler := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// If the context is missing the value, request the service to be shutdown gracefully
			v, ok := ctx.Value(web.KeyValues).(*web.Values)
			if !ok {
				return web.NewShutdownError("web value missing from context")
			}

			log.Printf(
				"%s : started   : %s %s -> %s",
				v.TraceId,
				r.Method, r.URL.Path, r.RemoteAddr,
			)

			err := innerHandler(ctx, w, r)

			log.Printf(
				"%s : completed : %s %s -> %s (%d) (%s)",
				v.TraceId,
				r.Method, r.URL.Path, r.RemoteAddr,
				v.StatusCode, time.Since(v.Now),
			)

			return err
		}

		return outerHandler
	}

	return middleware
}

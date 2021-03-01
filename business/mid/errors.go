package mid

import (
	"context"
	"github.com/jean-pasqualini/go-service/foundation/web"
	"log"
	"net/http"
)

// Errors ...
func Errors(log *log.Logger) web.Middleware {
	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// If the context is missing the value, request the service to be shutdown gracefully
			v, ok := ctx.Value(web.KeyValues).(*web.Values)
			if !ok {
				return web.NewShutdownError("web value missing from context")
			}

			if err := handler(ctx, w, r); err != nil {
				// Log the error
				log.Printf("%s : ERROR     : %v", v.TraceId, err)

				// Respond to the errors
				if err := web.RespondError(ctx, w, err); err != nil {
					return err
				}

				// If we receive the shutdown err we need to return it back to the base handler to shutdown the service.
				if ok := web.IsShutdown(err); ok {
					return err
				}
			}

			return nil
		}
	}
}

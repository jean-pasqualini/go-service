package mid

import (
	"context"
	"expvar"
	"github.com/jean-pasqualini/go-service/foundation/web"
	"net/http"
	"runtime"
)

// m contains the global program counters for the application.
var m = struct {
	gr  *expvar.Int
	req *expvar.Int
	err *expvar.Int
}{
	gr:  expvar.NewInt("goroutines"),
	req: expvar.NewInt("requests"),
	err: expvar.NewInt("errors"),
}

// Metrics updates program counters
// Factory a middleware
func Metrics() web.Middleware {

	// Decorator
	return func(handler web.Handler) web.Handler {

		// Decorate handler
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			err := handler(ctx, w, r)

			// Increment the request counter.
			m.req.Add(1)

			// Update the count for the number of active goroutines every 100 requests.
			// 50000 in production could be a good value
			if m.req.Value()%100 == 0 {
				m.gr.Set(int64(runtime.NumGoroutine()))
			}

			// Increment the errors counter if an error occurred on this request.
			if err != nil {
				m.err.Add(1)
			}

			// return the error so it can be handled further up the chain.
			return err
		}
	}
}

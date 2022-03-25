package mid

import (
	"context"
	"expvar"
	"net/http"
	"runtime"

	"github.com/naixyeur/garagesale/internal/platform/web"
)

// ================================================================================
var m = struct {
	gr  *expvar.Int
	req *expvar.Int
	err *expvar.Int
}{
	gr:  expvar.NewInt("goroutines"),
	req: expvar.NewInt("requests"),
	err: expvar.NewInt("errors"),
}

// ================================================================================
// Metrics
func Metrics() web.Middleware {

	return func(fn web.HandlerFunc) web.HandlerFunc {

		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			err := fn(ctx, w, r)

			m.req.Add(1)

			if m.req.Value()%100 == 0 {
				m.gr.Set(int64(runtime.NumGoroutine()))
			}

			if err != nil {
				m.err.Add(1)
			}

			return err

		}

	}

}

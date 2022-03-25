package mid

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/naixyeur/garagesale/internal/platform/web"
)

// ================================================================================
// Logger
func Logger(log *log.Logger) web.Middleware {

	return func(fn web.HandlerFunc) web.HandlerFunc {

		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			ctxVal, ok := ctx.Value(web.KeyValues).(*web.Values)
			if !ok {
				return errors.New("web values missing from context")
			}

			err := fn(ctx, w, r)

			fmt.Printf(
				"%v %s %s (%v)\n",
				ctxVal.StatusCode,
				r.Method,
				r.URL,
				time.Since(ctxVal.Start),
			)

			return err
		}

	}

}

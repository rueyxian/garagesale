package mid

import (
	"context"
	"log"
	"net/http"

	"github.com/naixyeur/garagesale/internal/platform/web"
)

// ================================================================================
// Errors
func Errors(log *log.Logger) web.Middleware {

	return func(fn web.HandlerFunc) web.HandlerFunc {

		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			if err := fn(ctx, w, r); err != nil {
				log.Printf("ERROR : %+v", err)
				return web.ResponseError(ctx, w, err)
			}

			return nil
		}

	}

}

package web

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

// ================================================================================
// ctxKey
type ctxKey int

const (
	KeyValues ctxKey = 1
)

// Values
type Values struct {
	StatusCode int
	Start      time.Time
}

// ================================================================================
// HandlerFunc
type HandlerFunc func(context.Context, http.ResponseWriter, *http.Request) error

// ================================================================================
// App is the entrypoint for all web application.
type App struct {
	mux *chi.Mux
	log *log.Logger
	mws []Middleware
}

// ================================================================================
// NewApp constructs internal state for an App.
func NewApp(log *log.Logger, mws ...Middleware) *App {
	return &App{
		mux: chi.NewRouter(),
		log: log,
		mws: mws,
	}
}

// ================================================================================
// Handle connects a method and URL pattern to a particular application handler.
func (a *App) Handle(method, pattern string, hfn HandlerFunc) {

	hfn = WrapMiddleware(a.mws, hfn)

	fn := func(w http.ResponseWriter, r *http.Request) {

		v := Values{
			Start: time.Now(),
		}

		ctx := context.WithValue(r.Context(), KeyValues, &v)

		if err := hfn(ctx, w, r); err != nil {
			a.log.Printf("unhandled error: %+v\n", err)
		}

	}
	a.mux.MethodFunc(method, pattern, fn)
}

// ================================================================================
// ServeHTTP
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

package handlers

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	mid "github.com/naixyeur/garagesale/internal/middleware"
	"github.com/naixyeur/garagesale/internal/platform/web"
)

func API(log *log.Logger, db *sqlx.DB) http.Handler {

	app := web.NewApp(log, mid.Logger(log), mid.Errors(log), mid.Metrics())

	{
		c := Check{DB: db}
		app.Handle(http.MethodGet, "/v1/health", c.Health)
	}

	p := Product{DB: db, Log: log}
	app.Handle(http.MethodGet, "/v1/products", p.List)
	app.Handle(http.MethodPost, "/v1/products", p.Create)
	app.Handle(http.MethodGet, "/v1/products/{id}", p.Retrieve)
	app.Handle(http.MethodPut, "/v1/products/{id}", p.Update)
	app.Handle(http.MethodDelete, "/v1/products/{id}", p.Delete)

	app.Handle(http.MethodGet, "/v1/products/{id}/sales", p.ListSale)
	app.Handle(http.MethodPost, "/v1/products/{id}/sales", p.AddSale)

	return app
}

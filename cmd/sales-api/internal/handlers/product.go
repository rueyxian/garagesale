package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/naixyeur/garagesale/internal/platform/web"
	"github.com/naixyeur/garagesale/internal/product"
	"github.com/pkg/errors"
)

type Product struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// ================================================================================
// List: web.HandlerFunc for listing all products
func (p *Product) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	list, err := product.List(ctx, p.DB)
	if err != nil {
		return err
	}

	return web.Response(ctx, w, list, http.StatusOK)
}

// ================================================================================
// Retrieve: web.HandlerFunc for retrieving a product
func (p *Product) Retrieve(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	prod, err := product.Retrieve(ctx, p.DB, id)
	if err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidId:
			return web.NewRequestError(err, http.StatusBadRequest)
		}
		return errors.Wrapf(err, "looking for product %q", id)
	}
	return web.Response(ctx, w, prod, http.StatusOK)
}

// ================================================================================
// Create: web.HandlerFunc for creating new Product
func (p *Product) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var np product.NewProduct
	if err := web.Decode(r, &np); err != nil {
		return errors.Wrap(err, "decoding client request body")
	}
	prod, err := product.Create(ctx, p.DB, np, time.Now())
	if err != nil {
		return err
	}
	return web.Response(ctx, w, prod, http.StatusOK)
}

// ================================================================================
// AddSale
func (p *Product) AddSale(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var ns product.NewSale

	if err := web.Decode(r, &ns); err != nil {
		return errors.Wrap(err, "decoding client request body")
	}

	productID := chi.URLParam(r, "id")

	sale, err := product.AddSale(ctx, p.DB, productID, ns, time.Now())
	if err != nil {
		return err
	}
	return web.Response(ctx, w, sale, http.StatusCreated)
}

// ================================================================================
// ListSale
func (p *Product) ListSale(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	list, err := product.ListSale(ctx, p.DB, id)
	if err != nil {
		return errors.Wrap(err, "getting sales list")
	}
	return web.Response(ctx, w, list, http.StatusOK)
}

// ================================================================================
// Update
func (p *Product) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	var up product.UpdateProduct
	if err := web.Decode(r, &up); err != nil {
		return errors.Wrap(err, "decoding client request body")
	}

	if err := product.Update(ctx, p.DB, id, up, time.Now()); err != nil {
		switch err {
		case product.ErrInvalidId:
			return web.NewRequestError(err, http.StatusBadRequest)
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		}
		return err
	}

	return web.Response(ctx, w, nil, http.StatusNoContent)
}

// ================================================================================
// Delete
func (p *Product) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	if err := product.Delete(ctx, p.DB, id); err != nil {
		switch err {
		case product.ErrInvalidId:
			return web.NewRequestError(err, http.StatusBadRequest)
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		}
		return err
	}

	return web.Response(ctx, w, nil, http.StatusNoContent)
}

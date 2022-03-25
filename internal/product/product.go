package product

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var (
	ErrNotFound  = errors.New("product not found")
	ErrInvalidId = errors.New("invalid UUID")
)

// ================================================================================
// List return a slice of Product
func List(ctx context.Context, db *sqlx.DB) ([]Product, error) {
	list := []Product{}
	// const q = `SELECT
	//   product_id, name, cost, quantity, date_created, date_updated
	//   FROM products`

	const q = `SELECT
		p.product_id, p.name, p.cost, p.quantity, 
		COALESCE(SUM(s.quantity), 0) AS sold,
		COALESCE(SUM(s.paid), 0) AS revenue,
		p.date_created, p.date_updated 
		FROM products AS p
		LEFT JOIN sales as s ON p.product_id = s.product_id
		GROUP BY p.product_id`

	if err := db.SelectContext(ctx, &list, q); err != nil {
		return nil, err
	}
	return list, nil
}

// ================================================================================
// Retrieve returns a single Product
func Retrieve(ctx context.Context, db *sqlx.DB, id string) (*Product, error) {

	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidId
	}

	var p Product
	const q = `SELECT
			p.product_id, p.name, p.cost, p.quantity, 
			COALESCE(SUM(s.quantity), 0) AS sold,
			COALESCE(SUM(s.paid), 0) AS revenue,
			p.date_created, p.date_updated 
		FROM products AS p
		LEFT JOIN sales as s ON p.product_id = s.product_id
		WHERE p.product_id = $1
		GROUP BY p.product_id`

	if err := db.GetContext(ctx, &p, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

// ================================================================================
// Create insert new Product
func Create(ctx context.Context, db *sqlx.DB, np NewProduct, now time.Time) (*Product, error) {
	p := Product{
		ID:          uuid.New().String(),
		Name:        np.Name,
		Cost:        np.Cost,
		Quantity:    np.Quantity,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}

	const q = `INSERT INTO products
	(product_id, name, cost, quantity, date_created, date_updated)
	VALUES($1, $2, $3, $4, $5, $6)`

	if _, err := db.ExecContext(ctx, q, p.ID, p.Name, p.Cost, p.Quantity, p.DateCreated, p.DateUpdated); err != nil {
		return nil, errors.Wrapf(err, "inserting product: %v", np)
	}
	return &p, nil
}

// ================================================================================
// Update
func Update(ctx context.Context, db *sqlx.DB, id string, up UpdateProduct, now time.Time) error {

	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidId
	}

	p, err := Retrieve(ctx, db, id)
	if err != nil {
		return nil
	}

	if up.Name != nil {
		p.Name = *up.Name
	}

	if up.Cost != nil {
		p.Cost = *up.Cost
	}

	if up.Quantity != nil {
		p.Quantity = *up.Quantity
	}

	p.DateUpdated = now

	const q = `UPDATE products SET
		"name" = $2,
		"cost" = $3,
		"quantity" = $4,
		"date_updated" = $5
		WHERE product_id = $1`

	if _, err := db.ExecContext(ctx, q, id, p.Name, p.Cost, p.Quantity, p.DateUpdated); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return errors.Wrap(err, "updating product")
	}

	return nil
}

// ================================================================================
// Delete
func Delete(ctx context.Context, db *sqlx.DB, id string) error {

	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidId
	}

	const q = `DELETE FROM products WHERE product_id = $1`

	if _, err := db.ExecContext(ctx, q, id); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return errors.Wrap(err, "deleting product")
	}

	return nil
}

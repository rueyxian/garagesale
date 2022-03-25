package product

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func AddSale(ctx context.Context, db *sqlx.DB, productID string, ns NewSale, now time.Time) (*Sale, error) {
	s := Sale{
		ID:          uuid.New().String(),
		ProductID:   productID,
		Quantity:    ns.Quantity,
		Paid:        ns.Paid,
		DateCreated: now.UTC(),
	}

	q := `INSERT INTO sales
		(sale_id, product_id, quantity, paid, date_created)
		VALUES($1, $2, $3, $4, $5)`

	if _, err := db.ExecContext(ctx, q, s.ID, s.ProductID, s.Quantity, s.Paid, s.DateCreated); err != nil {
		return nil, errors.Wrapf(err, "inserting sale: %v", ns)
	}
	return &s, nil
}

func ListSale(ctx context.Context, db *sqlx.DB, productID string) ([]Sale, error) {
	s := []Sale{}

	q := `SELECT
		sale_id, product_id, quantity, paid, date_created
		FROM sales
		WHERE product_id = $1`

	if err := db.SelectContext(ctx, &s, q, productID); err != nil {
		return nil, errors.Wrap(err, "selecting sales")
	}
	return s, nil
}

package product_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/naixyeur/garagesale/internal/product"
	"github.com/naixyeur/garagesale/internal/schema"
	"github.com/naixyeur/garagesale/internal/tests"
)

func TestProduct(t *testing.T) {
	db, cleanup := tests.NewUnit(t)
	defer cleanup()

	np := product.NewProduct{
		Name:     "Comic Book",
		Cost:     10,
		Quantity: 20,
	}

	now := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)

	p1, err := product.Create(db, np, now)
	if err != nil {
		t.Fatalf("could not create product: %v", err)
	}

	p2, err := product.Retrieve(db, p1.ID)
	if err != nil {
		t.Fatalf("could not retrieve product: %v", err)
	}

	if diff := cmp.Diff(p1, p2); diff != "" {
		t.Fatalf("fetch != created:\n %v", diff)
	}
}

func TestProductList(t *testing.T) {
	db, cleanup := tests.NewUnit(t)
	defer cleanup()

	if err := schema.Seed(db); err != nil {
		t.Fatal(err)
	}

	ps, err := product.List(db)
	if err != nil {
		t.Fatalf("listing prouct: %v", err)
	}

	if want, got := 2, len(ps); want != got {
		t.Fatalf("product size, want: %d, got: %d", want, got)
	}

}

package tests

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/naixyeur/garagesale/cmd/sales-api/internal/handlers"
	"github.com/naixyeur/garagesale/internal/schema"
	"github.com/naixyeur/garagesale/internal/tests"
)

func TestProducts(t *testing.T) {
	db, cleanup := tests.NewUnit(t)
	defer cleanup()

	if err := schema.Seed(db); err != nil {
		t.Fatal(err)
	}

	log := log.New(os.Stderr, "TEST : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	tests := ProductTests{app: handlers.API(log, db)}

	t.Run("List", tests.List)
	t.Run("ProductCRUD", tests.ProductCRUD)

}

type ProductTests struct {
	app http.Handler
}

func (p *ProductTests) List(t *testing.T) {
	req := httptest.NewRequest("GET", "/v1/products", nil)
	res := httptest.NewRecorder()

	p.app.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status code: want: %v, got: %v", http.StatusOK, res.Code)
	}

	var got []map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decoding: %s", err)
	}

	want := []map[string]interface{}{
		{
			"id":           "a2b0639f-2cc6-44b8-b97b-15d69dbb511e",
			"name":         "Comic Books",
			"cost":         float64(50),
			"quantity":     float64(42),
			"date_created": "2019-01-01T00:00:01.000001Z",
			"date_updated": "2019-01-01T00:00:01.000001Z",
		},
		{
			"id":           "72f8b983-3eb4-48db-9ed0-e45cc6bd716b",
			"name":         "McDonalds Toys",
			"cost":         float64(75),
			"quantity":     float64(120),
			"date_created": "2019-01-01T00:00:02.000001Z",
			"date_updated": "2019-01-01T00:00:02.000001Z",
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}

}

func (p *ProductTests) ProductCRUD(t *testing.T) {

	body := strings.NewReader(`{"name": "product0", "cost": 10, "quantity": 20}`)

	req := httptest.NewRequest("POST", "/v1/products", body)
	res := httptest.NewRecorder()

	p.app.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status code: want: %v, got: %v", http.StatusOK, res.Code)
	}

	var got map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decoding: %s", err)
	}

	if got["id"] == "" || got["id"] == nil {
		t.Fatal("expected non-empty product id")
	}

	if got["date_created"] == "" || got["date_created"] == nil {
		t.Fatal("expected non-empty product date_created")
	}

	if got["date_updated"] == "" || got["date_updated"] == nil {
		t.Fatal("expected non-empty product date_updated")
	}

	want := map[string]interface{}{
		"id":           got["id"],
		"name":         "product0",
		"cost":         float64(10),
		"quantity":     float64(20),
		"date_created": got["date_created"],
		"date_updated": got["date_updated"],
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}

}

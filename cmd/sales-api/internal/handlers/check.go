package handlers

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/naixyeur/garagesale/internal/platform/database"
	"github.com/naixyeur/garagesale/internal/platform/web"
)

// ================================================================================
// Check
type Check struct {
	DB *sqlx.DB
}

// ================================================================================
// Health
func (c *Check) Health(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	var health struct {
		Status string `json:"status"`
	}

	if err := database.StatusCheck(r.Context(), c.DB); err != nil {
		health.Status = "db not ready"
		return web.Response(ctx, w, health, http.StatusInternalServerError)
	}

	health.Status = "db ok"
	return web.Response(ctx, w, health, http.StatusOK)

}

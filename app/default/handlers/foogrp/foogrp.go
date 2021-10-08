// Package foogrp maintains the group of handlers for foo access.
package foogrp

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cpustejovsky/mongogo/internal/core/foo"

	v1Web "github.com/ardanlabs/service/business/web/v1"
	"github.com/ardanlabs/service/foundation/web"
)

// Handlers manages the set of foo enpoints.
type Handlers struct {
	Core foo.Core
}

// Create adds a new foo to the system.
func (h Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	var nu foo.Newfoo
	if err := web.Decode(r, &nu); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	usr, err := h.Core.Create(ctx, nu, v.Now)
	if err != nil {
		return fmt.Errorf("foo[%+v]: %w", &usr, err)
	}

	return web.Respond(ctx, w, usr, http.StatusCreated)
}

// Update updates a foo in the system.
func (h Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	var upd foo.Updatefoo
	if err := web.Decode(r, &upd); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	fooID := web.Param(r, "id")

	if err := h.Core.Update(ctx, fooID, upd, v.Now); err != nil {
		return fmt.Errorf("ID[%s] foo[%+v]: %w", fooID, &upd, err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Delete removes a foo from the system.
func (h Handlers) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	fooID := web.Param(r, "id")

	if err := h.Core.Delete(ctx, fooID); err != nil {
		return fmt.Errorf("ID[%s]: %w", fooID, err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Query returns a list of foos with paging.
func (h Handlers) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	page := web.Param(r, "page")
	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return v1Web.NewRequestError(fmt.Errorf("invalid page format [%s]", page), http.StatusBadRequest)
	}
	rows := web.Param(r, "rows")
	rowsPerPage, err := strconv.Atoi(rows)
	if err != nil {
		return v1Web.NewRequestError(fmt.Errorf("invalid rows format [%s]", rows), http.StatusBadRequest)
	}

	foos, err := h.Core.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		return fmt.Errorf("unable to query for foos: %w", err)
	}

	return web.Respond(ctx, w, foos, http.StatusOK)
}

// QueryByID returns a foo by its ID.
func (h Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	fooID := web.Param(r, "id")

	usr, err := h.Core.QueryByID(ctx, fooID)
	if err != nil {
		return fmt.Errorf("ID[%s]: %w", fooID, err)
	}

	return web.Respond(ctx, w, usr, http.StatusOK)
}

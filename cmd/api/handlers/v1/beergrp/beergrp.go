// Package beergrp maintains the group of handlers for beer access.
package beergrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/phbpx/gobeers/internal/core/beer"
	v1Web "github.com/phbpx/gobeers/internal/web/v1"
	"github.com/phbpx/gobeers/kit/web"
)

// Handlers manages the set of beer endpoints.
type Handlers struct {
	Beer beer.Core
}

// Create adds a new beer to the system.
func (h Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	var nb beer.NewBeer
	if err := web.Decode(r, &nb); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	b, err := h.Beer.Create(ctx, nb, v.Now)
	if err != nil {
		return fmt.Errorf("creating new beer, nb[%+v]: %w", nb, err)
	}

	return web.Respond(ctx, w, b, http.StatusCreated)
}

// QueryByID returns a beer by its ID.
func (h Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := web.Param(r, "id")
	b, err := h.Beer.QueryByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, beer.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, beer.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s]: %w", id, err)
		}
	}

	return web.Respond(ctx, w, b, http.StatusOK)
}

func (h Handlers) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	opts := v1Web.GetPagingOptions(r)

	list, err := h.Beer.Query(ctx, opts.Page, opts.PageSize)
	if err != nil {
		return fmt.Errorf("querying beers: %w", err)
	}

	return web.Respond(ctx, w, list, http.StatusOK)
}

// CreateReview adds a new review to an existing beer.
func (h Handlers) CreateReview(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	var nr beer.NewReview
	if err := web.Decode(r, &nr); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	id := web.Param(r, "id")

	rw, err := h.Beer.CreateReview(ctx, id, nr, v.Now)
	if err != nil {
		switch {
		case errors.Is(err, beer.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, beer.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("creating review ID[%s], nr[%+v]: %w", id, nr, err)
		}
	}

	return web.Respond(ctx, w, rw, http.StatusCreated)
}

// QueryReviews returns all reviews for a beer.
func (h Handlers) QueryReviews(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := web.Param(r, "id")
	opts := v1Web.GetPagingOptions(r)

	reviews, err := h.Beer.QueryReviews(ctx, id, opts.Page, opts.PageSize)
	if err != nil {
		switch {
		case errors.Is(err, beer.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, beer.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying reviews ID[%s]: %w", id, err)
		}
	}

	return web.Respond(ctx, w, reviews, http.StatusOK)
}
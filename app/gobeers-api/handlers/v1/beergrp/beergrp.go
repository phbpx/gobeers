// Package beergrp maintains the group of handlers for beer access.
package beergrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/phbpx/gobeers/business/core/beer"
	v1Web "github.com/phbpx/gobeers/business/web/v1"
	"github.com/phbpx/gobeers/foundation/web"
)

const (
	defaultPage = "1"
	defaultSize = "10"
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
	page := web.Query(r, "page", defaultPage)
	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return v1Web.NewRequestError(fmt.Errorf("invalid page format, page[%s]", page), http.StatusBadRequest)
	}

	size := web.Query(r, "size", defaultSize)
	sizeNumber, err := strconv.Atoi(size)
	if err != nil {
		return v1Web.NewRequestError(fmt.Errorf("invalid rows format, size[%s]", size), http.StatusBadRequest)
	}

	list, err := h.Beer.Query(ctx, pageNumber, sizeNumber)
	if err != nil {
		return fmt.Errorf("querying beers: %w", err)
	}

	if len(list) == 0 {
		return web.Respond(ctx, w, nil, http.StatusNoContent)
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

	page := web.Query(r, "page", defaultPage)
	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return v1Web.NewRequestError(fmt.Errorf("invalid page format, page[%s]", page), http.StatusBadRequest)
	}

	size := web.Query(r, "size", defaultSize)
	sizeNumber, err := strconv.Atoi(size)
	if err != nil {
		return v1Web.NewRequestError(fmt.Errorf("invalid rows format, size[%s]", size), http.StatusBadRequest)
	}

	reviews, err := h.Beer.QueryReviews(ctx, id, pageNumber, sizeNumber)
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

	if len(reviews) == 0 {
		return web.Respond(ctx, w, nil, http.StatusNoContent)
	}

	return web.Respond(ctx, w, reviews, http.StatusOK)
}

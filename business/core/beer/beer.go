// Package beer provides an example of a core business API. Right now these
// calls are just wrapping the data/store layer. But at some point you will
// want to audit or something that isn't specific to the data/store layer.
package beer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/phbpx/gobeers/business/sys/database"
	"github.com/phbpx/gobeers/business/sys/validate"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound  = errors.New("beer not found")
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// Storer .
type Storer interface {
	AddBeer(ctx context.Context, beer Beer) error
	QueryBeers(ctx context.Context, page int, size int) ([]Beer, error)
	QueryBeerByID(ctx context.Context, beerID string) (Beer, error)
	AddReview(ctx context.Context, review Review) error
	QueryBeerReviews(ctx context.Context, beerID string, page int, size int) ([]Review, error)
}

// Core manages the set of APIs for beer access.
type Core struct {
	store Storer
}

// NewCore constructs a core for product api access.
func NewCore(store Storer) Core {
	return Core{
		store: store,
	}
}

// =========================================================================
// Beer Support

// Create adds an beer to the database. Its return the created Beer
// with fields populated.
func (c Core) Create(ctx context.Context, nb NewBeer) (Beer, error) {
	if err := validate.Check(nb); err != nil {
		return Beer{}, fmt.Errorf("validating data: %w", err)
	}

	beer := Beer{
		ID:        uuid.New().String(),
		Name:      nb.Name,
		Brewery:   nb.Brewery,
		Style:     nb.Style,
		ABV:       nb.ABV,
		ShortDesc: nb.ShortDesc,
		CreatedAt: time.Now(),
	}

	if err := c.store.AddBeer(ctx, beer); err != nil {
		return Beer{}, fmt.Errorf("addBeer: %w", err)
	}

	return beer, nil
}

// QueryByID gets the specified beer from the database.
func (c Core) QueryByID(ctx context.Context, id string) (Beer, error) {
	if err := validate.CheckID(id); err != nil {
		return Beer{}, ErrInvalidID
	}

	beer, err := c.store.QueryBeerByID(ctx, id)
	if err != nil {
		if database.IsNoRowError(err) {
			return Beer{}, ErrNotFound
		}
		return Beer{}, fmt.Errorf("queryBeerByID: %w", err)
	}

	return beer, nil
}

// Query gets all beers from the database.
func (c Core) Query(ctx context.Context, page, pageSize int) ([]Beer, error) {
	beers, err := c.store.QueryBeers(ctx, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("queryBeers: %w", err)
	}

	return beers, nil
}

// =========================================================================
// Beer Review Support

// CreateReview adds a review to the database. Its return the created Review
// with fields populated.
func (c Core) CreateReview(ctx context.Context, beerID string, nr NewReview, now time.Time) (Review, error) {
	if err := validate.CheckID(beerID); err != nil {
		return Review{}, ErrInvalidID
	}

	if err := validate.Check(nr); err != nil {
		return Review{}, fmt.Errorf("validating data: %w", err)
	}

	beer, err := c.store.QueryBeerByID(ctx, beerID)
	if err != nil {
		if database.IsNoRowError(err) {
			return Review{}, ErrNotFound
		}
		return Review{}, fmt.Errorf("reviewing beer berrID[%s]: %w", beerID, err)
	}

	review := Review{
		ID:        uuid.New().String(),
		UserID:    nr.UserID,
		BeerID:    beer.ID,
		Score:     nr.Score,
		Comment:   nr.Comment,
		CreatedAt: now,
	}

	if err := c.store.AddReview(ctx, review); err != nil {
		return Review{}, fmt.Errorf("addReview: %w", err)
	}

	return review, nil
}

// QueryReviews gets all reviews for a beer from the database.
func (c Core) QueryReviews(ctx context.Context, beerID string, page int, size int) ([]Review, error) {
	if err := validate.CheckID(beerID); err != nil {
		return nil, ErrInvalidID
	}

	reviews, err := c.store.QueryBeerReviews(ctx, beerID, page, size)
	if err != nil {
		return nil, fmt.Errorf("queryBeerReviews: %w", err)
	}

	return reviews, nil
}

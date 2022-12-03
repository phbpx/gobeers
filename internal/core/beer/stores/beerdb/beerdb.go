// Package beerdb contains beer/review related CRUD functionality.
package beerdb

import (
	"context"
	"fmt"

	"github.com/phbpx/gobeers/internal/core/beer"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

// Store manages the set of APIs for beer access.
type Store struct {
	log *zap.SugaredLogger
	db  *bun.DB
}

// NewStore constructs a data for api access.
func NewStore(log *zap.SugaredLogger, db *bun.DB) Store {
	return Store{
		log: log,
		db:  db,
	}
}

// AddBeer adds a new beer to the database.
func (s Store) AddBeer(ctx context.Context, b beer.Beer) error {
	dbBeer := toDBBeer(b)

	if _, err := s.db.NewInsert().Model(&dbBeer).Exec(ctx); err != nil {
		return fmt.Errorf("adding beer: %w", err)
	}

	return nil
}

// QueryBeerByID retrieves a beer by its id.
func (s Store) QueryBeerByID(ctx context.Context, beerID string) (beer.Beer, error) {
	var b dbBeer

	query := s.db.NewSelect().
		Model(&b).
		Where("id = ?", beerID)

	if err := query.Scan(ctx); err != nil {
		return beer.Beer{}, fmt.Errorf("querying beer by [id=%s]: %w", beerID, err)
	}

	return toBeer(b), nil
}

// QueryBeers retrieves a list of existing beers.
func (s Store) QueryBeers(ctx context.Context, page, size int) ([]beer.Beer, error) {
	var beers []dbBeer

	query := s.db.NewSelect().
		Model(&beers).
		Limit(size).
		Offset(size * (page - 1))

	if err := query.Scan(ctx); err != nil {
		return nil, fmt.Errorf("querying beer: %w", err)
	}

	return toBeers(beers), nil
}

// AddReview adds a new beer review to the database.
func (s Store) AddReview(ctx context.Context, r beer.Review) error {
	dbReview := toDBReview(r)

	if _, err := s.db.NewInsert().Model(&dbReview).Exec(ctx); err != nil {
		return fmt.Errorf("adding review: %w", err)
	}

	return nil
}

// QueryBeerReviews retrieves a list of reviews for a beer.
func (s Store) QueryBeerReviews(ctx context.Context, beerID string, page int, size int) ([]beer.Review, error) {
	var reviews []dbReview

	query := s.db.NewSelect().
		Model(&reviews).
		Where("beer_id =?", beerID).
		Limit(size).
		Offset(size * (page - 1))

	if err := query.Scan(ctx); err != nil {
		return nil, fmt.Errorf("querying beer review [beer_id=%s]: %w", beerID, err)
	}

	return toReviews(reviews), nil
}

package beer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/phbpx/gobeers/internal/core/beer/db"
	"github.com/phbpx/gobeers/internal/sys/database"
	"github.com/phbpx/gobeers/internal/sys/validate"
	"go.uber.org/zap"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound  = errors.New("beer not found")
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// Core manages the set of APIs for beer access.
type Core struct {
	store db.Store
}

// NewCore constructs a core for product api access.
func NewCore(log *zap.SugaredLogger, sqlxDB *sqlx.DB) Core {
	return Core{
		store: db.NewStore(log, sqlxDB),
	}
}

// =========================================================================
// Beer Support

// Create adds an beer to the database. Its return the created Beer
// with fields populated.
func (c Core) Create(ctx context.Context, nb NewBeer, now time.Time) (Beer, error) {
	if err := validate.Check(nb); err != nil {
		return Beer{}, fmt.Errorf("validating data: %w", err)
	}

	dbBeer := db.Beer{
		ID:        uuid.New().String(),
		Name:      nb.Name,
		Brewery:   nb.Brewery,
		Style:     nb.Style,
		ABV:       nb.ABV,
		ShortDesc: nb.ShortDesc,
		CreatedAt: now,
	}

	if err := c.store.AddBeer(ctx, dbBeer); err != nil {
		return Beer{}, fmt.Errorf("addBeer: %w", err)
	}

	return toBeer(dbBeer), nil
}

// QueryByID gets the specified beer from the database.
func (c Core) QueryByID(ctx context.Context, id string) (Beer, error) {
	if err := validate.CheckID(id); err != nil {
		return Beer{}, ErrInvalidID
	}

	dbBeer, err := c.store.QueryBeerByID(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Beer{}, ErrNotFound
		}
		return Beer{}, fmt.Errorf("queryBeerByID: %w", err)
	}

	return toBeer(dbBeer), nil
}

// Query gets all beers from the database.
func (c Core) Query(ctx context.Context, page, pageSize int) ([]Beer, error) {
	dbBeers, err := c.store.QueryBeers(ctx, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("queryBeers: %w", err)
	}

	return toBeerSlice(dbBeers), nil
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

	_, err := c.store.QueryBeerByID(ctx, beerID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Review{}, ErrNotFound
		}
		return Review{}, fmt.Errorf("reviewing beer berrID[%s]: %w", beerID, err)
	}

	dbReview := db.Review{
		ID:        uuid.New().String(),
		UserID:    nr.UserID,
		BeerID:    beerID,
		Score:     nr.Score,
		Comment:   nr.Comment,
		CreatedAt: now,
	}

	if err := c.store.AddReview(ctx, dbReview); err != nil {
		return Review{}, fmt.Errorf("addReview: %w", err)
	}

	return toReview(dbReview), nil
}

// QueryReviews gets all reviews for a beer from the database.
func (c Core) QueryReviews(ctx context.Context, beerID string, page, pageSize int) ([]Review, error) {
	if err := validate.CheckID(beerID); err != nil {
		return nil, ErrInvalidID
	}

	dbReviews, err := c.store.QueryBeerReviews(ctx, beerID, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("queryBeerReviews: %w", err)
	}

	return toReviewSlice(dbReviews), nil
}

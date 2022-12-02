// Package beerdb contains beer/review related CRUD functionality.
package beerdb

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/phbpx/gobeers/internal/core/beer"
	"github.com/phbpx/gobeers/internal/sys/database"
	"go.uber.org/zap"
)

// Store manages the set of APIs for beer access.
type Store struct {
	log          *zap.SugaredLogger
	db           sqlx.ExtContext
}

// NewStore constructs a data for api access.
func NewStore(log *zap.SugaredLogger, db *sqlx.DB) Store {
	return Store{
		log: log,
		db:  db,
	}
}

// AddBeer adds a new beer to the database.
func (s Store) AddBeer(ctx context.Context, b beer.Beer) error {
    dbBeer := toDBBeer(b)

	const q = `
    INSERT INTO beers (
            id, 
            name, 
            brewery, 
            style, 
            abv, 
            short_desc, 
            created_at
    ) VALUES (
            :id, :name, :brewery, :style, :abv, :short_desc, :created_at
    )`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, &dbBeer); err != nil {
		return fmt.Errorf("inserting beer: %w", err)
	}

	return nil
}

// QueryBeerByID retrieves a beer by its id.
func (s Store) QueryBeerByID(ctx context.Context, beerID string) (beer.Beer, error) {
	params := struct {
		ID string `db:"id"`
	}{
		ID: beerID,
	}

	const q = `
    SELECT 
            b.id,
            b.name,
            b.brewery,
            b.style,
            b.abv,
            b.short_desc,
            COALESCE(AVG(r.score), 0) AS score,
            b.created_at
    FROM 
            beers AS b
    LEFT JOIN 
            reviews AS r ON r.beer_id = b.id
    WHERE 
            b.id = :id
    GROUP BY
            b.id`

	var b dbBeer
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, params, &b); err != nil {
		return beer.Beer{}, fmt.Errorf("selecting beer id[%q]: %w", beerID, err)
	}

	return toBeer(b), nil
}

// QueryBeers retrieves a list of existing beers.
func (s Store) QueryBeers(ctx context.Context, page, size int) ([]beer.Beer, error) {
	params := struct {
		Offset   int `db:"offset"`
		PageSize int `db:"page_size"`
	}{
		Offset:   (page - 1) * size,
		PageSize: size,
	}

	const q = `
    SELECT 
            b.id,
            b.name,
            b.brewery,
            b.style,
            b.abv,
            b.short_desc,
            COALESCE(AVG(r.score), 0) AS score,
            b.created_at
    FROM 
            beers AS b
    LEFT JOIN 
            reviews AS r ON r.beer_id = b.id
    GROUP BY
            b.id
    OFFSET :offset ROWS FETCH NEXT :page_size ROWS ONLY`

	var beers []dbBeer
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, params, &beers); err != nil {
		return []beer.Beer{}, fmt.Errorf("selecting beers: %w", err)
	}

	return toBeers(beers), nil
}

// AddReview adds a new beer review to the database.
func (s Store) AddReview(ctx context.Context, r beer.Review) error {
    dbReview := toDBReview(r)

	const q = `
    INSERT INTO reviews (
        id, 
        user_id,
        beer_id, 
        score, 
        comment, 
        created_at
    ) VALUES (
        :id, :user_id, :beer_id, :score, :comment, :created_at
    )`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, &dbReview); err != nil {
		return fmt.Errorf("inserting review: %w", err)
	}

	return nil
}

// QueryBeerReviews retrieves a list of reviews for a beer.
func (s Store) QueryBeerReviews(ctx context.Context, beerID string, page int, size int) ([]beer.Review, error) {
	params := struct {
		BeerID   string `db:"beer_id"`
		Offset   int    `db:"offset"`
		PageSize int    `db:"page_size"`
	}{
		BeerID:   beerID,
		Offset:   (page - 1) * size,
		PageSize: size,
	}

	const q = `
    SELECT 
            r.id,
            r.user_id,
            r.beer_id,
            r.score,
            r.comment,
            r.created_at
    FROM 
            reviews AS r
    WHERE 
            r.beer_id = :beer_id
    ORDER BY 
            r.created_at
    OFFSET :offset ROWS FETCH NEXT :page_size ROWS ONLY`

	var reviews []dbReview
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, params, &reviews); err != nil {
		return nil, fmt.Errorf("querying reviews: %w", err)
	}

	return toReviews(reviews), nil
}

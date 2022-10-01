package db

import (
	"context"
	"fmt"

	"github.com/phbpx/gobeers/internal/sys/database"
)

// AddReview adds a new beer review to the database.
func (s Store) AddReview(ctx context.Context, review Review) error {
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

	if err := database.NamedExecContext(ctx, s.log, s.db, q, review); err != nil {
		return fmt.Errorf("inserting review: %w", err)
	}

	return nil
}

// QueryBeerReviews retrieves a list of reviews for a beer.
func (s Store) QueryBeerReviews(ctx context.Context, beerID string, page, pageSize int) ([]Review, error) {
	params := struct {
		BeerID   string `db:"beer_id"`
		Offset   int    `db:"offset"`
		PageSize int    `db:"page_size"`
	}{
		BeerID:   beerID,
		Offset:   (page - 1) * pageSize,
		PageSize: pageSize,
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

	var reviews []Review
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, params, &reviews); err != nil {
		return nil, fmt.Errorf("querying reviews: %w", err)
	}

	return reviews, nil
}

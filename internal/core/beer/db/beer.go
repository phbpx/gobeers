package db

import (
	"context"
	"fmt"

	"github.com/phbpx/gobeers/internal/sys/database"
)

// AddBeer adds a new beer to the database.
func (s Store) AddBeer(ctx context.Context, beer Beer) error {
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

	if err := database.NamedExecContext(ctx, s.log, s.db, q, beer); err != nil {
		return fmt.Errorf("inserting beer: %w", err)
	}

	return nil
}

// QueryBeerByID retrieves a beer by its id.
func (s Store) QueryBeerByID(ctx context.Context, id string) (Beer, error) {
	params := struct {
		ID string `db:"id"`
	}{
		ID: id,
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

	var b Beer
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, params, &b); err != nil {
		return Beer{}, fmt.Errorf("selecting beer id[%q]: %w", id, err)
	}

	return b, nil
}

// QueryBeers retrieves a list of existing beers.
func (s Store) QueryBeers(ctx context.Context, page, pageSize int) ([]Beer, error) {
	params := struct {
		Offset   int `db:"offset"`
		PageSize int `db:"page_size"`
	}{
		Offset:   (page - 1) * pageSize,
		PageSize: pageSize,
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

	var beers []Beer
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, params, &beers); err != nil {
		return []Beer{}, fmt.Errorf("selecting beers: %w", err)
	}

	return beers, nil
}

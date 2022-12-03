package beerdb

import (
	"time"

	"github.com/phbpx/gobeers/internal/core/beer"
	"github.com/uptrace/bun"
)

// dbBeer represents an individual beer.
type dbBeer struct {
	bun.BaseModel `bun:"table:beers,alias:b"`

	ID        string    `bun:"id,pk"`
	Name      string    `bun:"name"`
	Brewery   string    `bun:"brewery"`
	Style     string    `bun:"style"`
	ABV       float32   `bun:"abv"`
	ShortDesc string    `bun:"short_desc"`
	CreatedAt time.Time `bun:"created_at"`
}

// dbReview defines the properties of a review.
type dbReview struct {
	bun.BaseModel `bun:"table:reviews,alias:r"`

	ID        string    `bun:"id,pk"`
	BeerID    string    `bun:"beer_id"`
	UserID    string    `bun:"user_id"`
	Score     float32   `bun:"score"`
	Comment   string    `bun:"comment"`
	CreatedAt time.Time `bun:"created_at"`
}

// =========================================================

func toDBBeer(b beer.Beer) dbBeer {
	return dbBeer{
		ID:        b.ID,
		Name:      b.Name,
		Brewery:   b.Brewery,
		Style:     b.Style,
		ABV:       b.ABV,
		ShortDesc: b.ShortDesc,
		CreatedAt: b.CreatedAt,
	}
}

func toBeer(b dbBeer) beer.Beer {
	return beer.Beer{
		ID:        b.ID,
		Name:      b.Name,
		Brewery:   b.Brewery,
		Style:     b.Style,
		ABV:       b.ABV,
		ShortDesc: b.ShortDesc,
		CreatedAt: b.CreatedAt,
	}
}

func toBeers(list []dbBeer) []beer.Beer {
	beers := make([]beer.Beer, len(list))
	for i, b := range list {
		beers[i] = toBeer(b)
	}
	return beers
}

func toDBReview(r beer.Review) dbReview {
	return dbReview{
		ID:        r.ID,
		BeerID:    r.BeerID,
		UserID:    r.UserID,
		Score:     r.Score,
		Comment:   r.Comment,
		CreatedAt: r.CreatedAt,
	}
}

func toReview(r dbReview) beer.Review {
	return beer.Review{
		ID:        r.ID,
		BeerID:    r.BeerID,
		UserID:    r.UserID,
		Score:     r.Score,
		Comment:   r.Comment,
		CreatedAt: r.CreatedAt,
	}
}

func toReviews(list []dbReview) []beer.Review {
	reviews := make([]beer.Review, len(list))
	for i, r := range list {
		reviews[i] = toReview(r)
	}
	return reviews
}

package beerdb

import (
	"time"

	"github.com/phbpx/gobeers/internal/core/beer"
)

// dbBeer represents an individual beer.
type dbBeer struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Brewery   string    `db:"brewery"`
	Style     string    `db:"style"`
	ABV       float32   `db:"abv"`
	ShortDesc string    `db:"short_desc"`
	Score     float32   `db:"score"`
	CreatedAt time.Time `db:"created_at"`
}

// dbReview defines the properties of a review.
type dbReview struct {
	ID        string    `db:"id"`
	BeerID    string    `db:"beer_id"`
	UserID    string    `db:"user_id"`
	Score     float32   `db:"score"`
	Comment   string    `db:"comment"`
	CreatedAt time.Time `db:"created_at"`
}

// =========================================================

func toDBBeer(b beer.Beer) dbBeer {
    return dbBeer(b)
}

func toBeer(b dbBeer) beer.Beer {
    return beer.Beer(b)
}

func toBeers(list []dbBeer) []beer.Beer {
    beers := make([]beer.Beer, len(list))
    for i, b := range list {
        beers[i] = toBeer(b)
    }
    return beers
}

func toDBReview(r beer.Review) dbReview {
    return dbReview(r)
}

func toReview(r dbReview) beer.Review {
    return beer.Review(r)
}

func toReviews(list []dbReview) []beer.Review {
    reviews := make([]beer.Review, len(list))
    for i, r := range list {
        reviews[i] = toReview(r)
    }
    return reviews
}

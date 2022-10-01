package beer

import (
	"time"

	"github.com/phbpx/gobeers/internal/core/beer/db"
)

// NewBeer represents a new beer to be added to the system.
type NewBeer struct {
	Name      string  `json:"name" validate:"required"`
	Brewery   string  `json:"brewery" validate:"required"`
	Style     string  `json:"style" validate:"required"`
	ABV       float32 `json:"abv" validate:"required"`
	ShortDesc string  `json:"short_desc" validate:"required"`
}

// Beer defines the properties of a beer.
type Beer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Brewery   string    `json:"brewery"`
	Style     string    `json:"style"`
	ABV       float32   `json:"abv"`
	ShortDesc string    `json:"short_desc"`
	Score     float32   `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}

// NewReview defines the input parameters for creating a new review.
type NewReview struct {
	UserID  string  `json:"user_id" validate:"required,uuid"`
	Score   float32 `json:"score" validate:"required"`
	Comment string  `json:"comment" validate:"required"`
}

// Review defines the properties of a review.
type Review struct {
	ID        string    `json:"id"`
	BeerID    string    `json:"beer_id"`
	UserID    string    `json:"user_id"`
	Score     float32   `json:"score"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

// =============================================================================

func toBeer(dbBeer db.Beer) Beer {
	return Beer{
		ID:        dbBeer.ID,
		Name:      dbBeer.Name,
		Brewery:   dbBeer.Brewery,
		Style:     dbBeer.Style,
		ABV:       dbBeer.ABV,
		ShortDesc: dbBeer.ShortDesc,
		Score:     dbBeer.Score,
		CreatedAt: dbBeer.CreatedAt,
	}
}

func toBeerSlice(dbBeers []db.Beer) []Beer {
	beers := make([]Beer, len(dbBeers))
	for i := range dbBeers {
		beers[i] = toBeer(dbBeers[i])
	}
	return beers
}

func toReview(dbReview db.Review) Review {
	return Review{
		ID:        dbReview.ID,
		BeerID:    dbReview.BeerID,
		UserID:    dbReview.UserID,
		Score:     dbReview.Score,
		Comment:   dbReview.Comment,
		CreatedAt: dbReview.CreatedAt,
	}
}

func toReviewSlice(dbReviews []db.Review) []Review {
	reviews := make([]Review, len(dbReviews))
	for i := range dbReviews {
		reviews[i] = toReview(dbReviews[i])
	}
	return reviews
}

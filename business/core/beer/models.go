package beer

import "time"

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

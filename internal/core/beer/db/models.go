package db

import "time"

// Beer represents an individual beer.
type Beer struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Brewery   string    `db:"brewery"`
	Style     string    `db:"style"`
	ABV       float32   `db:"abv"`
	ShortDesc string    `db:"short_desc"`
	Score     float32   `db:"score"`
	CreatedAt time.Time `db:"created_at"`
}

// Review defines the properties of a review.
type Review struct {
	ID        string    `db:"id"`
	BeerID    string    `db:"beer_id"`
	UserID    string    `db:"user_id"`
	Score     float32   `db:"score"`
	Comment   string    `db:"comment"`
	CreatedAt time.Time `db:"created_at"`
}

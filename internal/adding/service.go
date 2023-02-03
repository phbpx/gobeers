// Package adding defines the use case for adding a beer.
package adding

import (
	"context"
	"time"
	"strings"

	"github.com/google/uuid"
	"github.com/phbpx/gobeer/internal/beers"
	"github.com/phbpx/gobeer/internal/coffees"
)

// NewBeer represents a new beer to be added to the system.
type NewBeer struct {
	Name      string  `json:"name" binding:"required"`
	Brewery   string  `json:"brewery" binding:"required"`
	Style     string  `json:"style" binding:"required"`
	ABV       float32 `json:"abv" binding:"required"`
	ShortDesc string  `json:"short_desc" binding:"required"`
}

// NewBeer represents a new coffee to be added to the system.
type NewCoffee struct {
	Name      	string  `json:"name" binding:"required"`
	State     	string  `json:"state" binding:"required"`
	Bitterness	string  `json:"bitterness" binding:"required"`
	Acidity    	float32 `json:"acidity" binding:"required"`
	ShortDesc 	string  `json:"short_desc" binding:"required"`
}

// Repository defines the interface for the adding service to interact
// with the storage.
type Repository interface {
	// CreateBeer adds a new beer to the storage.
	CreateBeer(ctx context.Context, b beers.Beer) error
	// BeerExists checks if a beer with the given name and brewery already exists.
	BeerExists(ctx context.Context, name, brewery string) (bool, error)
	//CreateCoffee adds a new cofee to the storage.
	CreateCoffee(ctx context.Context, coffee coffees.Coffee) error
}

// Service provides adding operations.
type Service struct {
	r Repository
}

// NewService creates an adding service with the necessary dependencies.
func NewService(r Repository) *Service {
	return &Service{r}
}

// AddBeer adds a new beer to the system.
func (s *Service) AddBeer(ctx context.Context, b NewBeer) (*beers.Beer, error) {
	beer := beers.Beer{
		ID:        uuid.NewString(),
		Name:      b.Name,
		Brewery:   b.Brewery,
		Style:     b.Style,
		ABV:       b.ABV,
		ShortDesc: b.ShortDesc,
		Score:     0,
		CreatedAt: time.Now(),
	}

	// Check if the beer already exists.
	exists, err := s.r.BeerExists(ctx, beer.Name, beer.Brewery)
	if err != nil {
		return nil, err
	}

	// If the beer already exists, return an error.
	if exists {
		return nil, beers.ErrAlreadyExists
	}

	return &beer, s.r.CreateBeer(ctx, beer)
}

// AddCoffee adds a new coffee to the system.
func (s *Service) AddCoffee(ctx context.Context, c NewCoffee) (*coffees.Coffee, error) {
	coffee := coffees.Coffee{
		ID:        		uuid.NewString(),
		Name:			c.Name,
		State:			strings.ToUpper(c.State),
		Bitterness:     strings.ToUpper(c.Bitterness),
		Acidity:       	c.Acidity,
		ShortDesc: 		c.ShortDesc,
		CreatedAt: 		time.Now(),
	}

	//Check if state is valid
	if !coffees.ValidCoffeeState(coffee.State) {
		return nil, coffees.ErrInvalidCoffeeState
	}

	//Check if bitterness is valid
	if !coffees.ValidCoffeeBitterness(coffee.Bitterness) {
		return nil, coffees.ErrInvalidCoffeeBitterness
	}

	// // Check if the beer already exists.
	// exists, err := s.r.BeerExists(ctx, beer.Name, beer.Brewery)
	// if err != nil {
	// 	return nil, err
	// }

	// // If the beer already exists, return an error.
	// if exists {
	// 	return nil, beers.ErrAlreadyExists
	// }

	return &coffee, s.r.CreateCoffee(ctx, coffee)
}

package coffees

import (
	"errors"
	"time"
	"strings"
)

var (
	ErrInvalidID = errors.New("invalid coffee ID")
	ErrNotFound = errors.New("coffee not found")
	ErrAlreadyExists = errors.New("coffee already exists")
	ErrInvalidCoffeeState = errors.New("invalid coffee state")
	ErrInvalidCoffeeBitterness = errors.New("invalid coffee bitterness")
)

type Coffee struct {
	ID			string		`json:"id"`
	Name		string		`json:"name"`
	State		string		`json:"state"` // ["SP", "MG", "BA"]
	Bitterness	string		`json:"bitterness"` // ["low", "medium", "high"]
	Acidity		float32		`json:"acidity"`
	ShortDesc	string		`json:"short_desc"`
	CreatedAt	time.Time	`json:"created_at"`
}

var validStates = []string{"SP", "MG", "BH"}
var validBitterness = []string{"LOW", "MEDIUM", "HIGH"}

func Contains (s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func ValidCoffeeState(state string) bool {
	uState := strings.ToUpper(state)
	return Contains(validStates, uState)
}

func ValidCoffeeBitterness(bitterness string) bool {
	uBitterness := strings.ToUpper(bitterness)
	return Contains(validBitterness, uBitterness)
}

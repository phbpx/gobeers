// Package v1 contains the full set of handler functions and routes
// supported by the v1 web api.
package v1

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/phbpx/gobeers/cmd/api/handlers/v1/beergrp"
	"github.com/phbpx/gobeers/internal/core/beer"
	"github.com/phbpx/gobeers/kit/web"
	"go.uber.org/zap"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log *zap.SugaredLogger
	DB  *sqlx.DB
}

// Routes binds all the version 1 routes.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	// Register beer endpoints.
	bgh := beergrp.Handlers{
		Beer: beer.NewCore(cfg.Log, cfg.DB),
	}
	app.Handle(http.MethodGet, version, "/beers", bgh.Query)
	app.Handle(http.MethodGet, version, "/beers/:id", bgh.QueryByID)
	app.Handle(http.MethodPost, version, "/beers", bgh.Create)
	app.Handle(http.MethodPost, version, "/beers/:id", bgh.CreateReview)
	app.Handle(http.MethodPost, version, "/beers/:id/reviews", bgh.QueryReviews)
}

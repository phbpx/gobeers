package beer_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/phbpx/gobeers/business/core/beer"
	"github.com/phbpx/gobeers/business/core/beer/stores/beerdb"
	"github.com/phbpx/gobeers/business/data/dbtest"
	"github.com/phbpx/gobeers/foundation/docker"
)

var c *docker.Container

func TestMain(m *testing.M) {
	var err error
	c, err = dbtest.StartDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbtest.StopDB(c)

	m.Run()
}

func TestBeer(t *testing.T) {
	log, db, teardown := dbtest.NewUnit(t, c, "testbeer")
	t.Cleanup(teardown)

	core := beer.NewCore(beerdb.NewStore(log, db))

	t.Log("Given the need to work with Beer records.")
	{
		t.Logf("\tWhen handling a single Beer.")
		{
			ctx := context.Background()
			now := time.Date(2022, 2, 25, 0, 0, 0, 0, time.UTC)

			nb := beer.NewBeer{
				Name:      "Test Beer",
				Brewery:   "Test Brewery",
				Style:     "Test Style",
				ABV:       5.5,
				ShortDesc: "Test Short Description",
			}

			beer, err := core.Create(ctx, nb, now)
			if err != nil {
				t.Fatalf("\t [ERROR] Should be able to add a beer : %s", err)
			}
			t.Logf("\t [SUCCESS] Should be able to add a beer.")

			saved, err := core.QueryByID(ctx, beer.ID)
			if err != nil {
				t.Fatalf("\t [ERROR] Should be able to query a beer by id : %s", err)
			}
			t.Logf("\t [SUCCESS] Should be able to query a beer by id.")

			if diff := cmp.Diff(beer, saved); diff != "" {
				t.Fatalf("\t [ERROR] Should get back the same beer : %s", diff)
			}
			t.Logf("\t [SUCCESS] Should get back the same beer.")

			beers, err := core.Query(ctx, 1, 10)
			if err != nil {
				t.Fatalf("\t [ERROR] Should be able to query beers : %s", err)
			}
			t.Logf("\t [SUCCESS] Should be able to query beers.")

			if len(beers) == 0 {
				t.Fatalf("\t [ERROR] Should get back at least one beer.")
			}
			t.Logf("\t [SUCCESS] Should get back at least one beer.")
		}
	}

	t.Log("Given the need to work with Beer Review records.")
	{
		t.Logf("\tWhen handling a single Beer Review.")
		{
			ctx := context.Background()
			now := time.Date(2022, 2, 25, 0, 0, 0, 0, time.UTC)

			nb := beer.NewBeer{
				Name:      "Test Beer",
				Brewery:   "Test Brewery",
				Style:     "Test Style",
				ABV:       5.5,
				ShortDesc: "Test Short Description",
			}

			b, err := core.Create(ctx, nb, now)
			if err != nil {
				t.Fatalf("\t [ERROR] Should be able to add a beer : %s", err)
			}

			nr := beer.NewReview{
				UserID:  uuid.NewString(),
				Score:   5,
				Comment: "Test Comment",
			}

			_, err = core.CreateReview(ctx, b.ID, nr, now)
			if err != nil {
				t.Fatalf("\t [ERROR] Should be able to add a review : %s", err)
			}
			t.Logf("\t [SUCCESS] Should be able to add a review.")

			reviews, err := core.QueryReviews(ctx, b.ID, 1, 10)
			if err != nil {
				t.Fatalf("\t [ERROR] Should be able to query reviews : %s", err)
			}
			t.Logf("\t [SUCCESS] Should be able to query reviews.")

			if len(reviews) == 0 {
				t.Fatalf("\t [ERROR] Should get back at least one review.")
			}
			t.Logf("\t [SUCCESS] Should get back at least one review.")
		}
	}
}

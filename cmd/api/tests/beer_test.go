package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	v1Web "github.com/ardanlabs/service/business/web/v1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/phbpx/gobeers/cmd/api/handlers"
	"github.com/phbpx/gobeers/internal/data/dbtest"
	"github.com/phbpx/gobeers/internal/sys/validate"
)

// BeerTests holds methods for each beer subtest. This type allows
// passing dependencies for tests while still providing a convenient syntax
// when subtests are registered.
type BeerTests struct {
	app http.Handler
}

func TestBeers(t *testing.T) {
	t.Parallel()

	test := dbtest.NewIntegration(t, c, "inttestprods")
	t.Cleanup(test.Teardown)

	shutdown := make(chan os.Signal, 1)
	tests := BeerTests{
		app: handlers.APIMux(handlers.APIMuxConfig{
			Shutdown: shutdown,
			Log:      test.Log,
			DB:       test.DB,
		}),
	}

	t.Run("postBeers400", tests.postBeers400)
	t.Run("getBeers400", tests.getBeers400)
	t.Run("getBeers404", tests.getBeers404)
}

// postBeers400 validates a beer can't be created with the endpoint
// unless a valid beer payload is submitted.
func (bt *BeerTests) postBeers400(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/v1/beers", strings.NewReader(`{}`))
	w := httptest.NewRecorder()

	bt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate a new beer can't be created with an invalid payload.")
	{
		t.Log("\t When using an incomplete beer value.")
		{
			if w.Code != http.StatusBadRequest {
				t.Fatalf("\t [ERROR] Should receive a status code of 400 for the response : %v", w.Code)
			}
			t.Log("\t [SUCCESS] Should receive a status code of 400 for the response.")

			// Inspect the response.
			var got v1Web.ErrorResponse
			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
				t.Fatalf("\t [ERROR] Should be able to unmarshal the response to an error type : %v", err)
			}
			t.Log("\t [SUCCESS] Should be able to unmarshal the response to an error type.")

			fields := validate.FieldErrors{
				{Field: "name", Error: "name is a required field"},
				{Field: "brewery", Error: "brewery is a required field"},
				{Field: "style", Error: "style is a required field"},
				{Field: "abv", Error: "abv is a required field"},
				{Field: "short_desc", Error: "short_desc is a required field"},
			}
			exp := v1Web.ErrorResponse{
				Error:  "data validation error",
				Fields: fields.Fields(),
			}

			// We can't rely on the order of the field errors so they have to be
			// sorted. Tell the cmp package how to sort them.
			sorter := cmpopts.SortSlices(func(a, b validate.FieldError) bool {
				return a.Field < b.Field
			})

			if diff := cmp.Diff(got, exp, sorter); diff != "" {
				t.Fatalf("\t [ERROR] Should get the expected result. Diff:\n%s", diff)
			}
			t.Log("\t [SUCCESS] Should get the expected result.")
		}
	}
}

// getBeers400 validates a beer request for a malformed id.
func (bt *BeerTests) getBeers400(t *testing.T) {
	id := "12345"

	r := httptest.NewRequest(http.MethodGet, "/v1/beers/"+id, nil)
	w := httptest.NewRecorder()

	bt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate getting a product with a malformed id.")
	{
		t.Logf("\t When using the new beer %s.", id)
		{
			if w.Code != http.StatusBadRequest {
				t.Fatalf("\t [ERROR] Should receive a status code of 400 for the response : %v", w.Code)
			}
			t.Log("\t [SUCCESS] Should receive a status code of 400 for the response.")

			got := w.Body.String()
			exp := `{"error":"ID is not in its proper form"}`
			if got != exp {
				t.Fatalf("\t [ERROR] Should get the expected result.\n\t\t Got: %s.\n\t\t Exp: %s", got, exp)
			}
			t.Log("\t [SUCCESS] Should get the expected result.")
		}
	}
}

func (bt *BeerTests) getBeers404(t *testing.T) {
	id := "112262f1-1a77-4374-9f22-39e575aa6348"

	r := httptest.NewRequest(http.MethodGet, "/v1/beers/"+id, nil)
	w := httptest.NewRecorder()

	bt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate deleting a product that does not exist.")
	{
		t.Log("\t Given the need to validate getting a beer with an unknown id.")
		{
			if w.Code != http.StatusNotFound {
				t.Fatalf("\t [ERROR] Should receive a status code of 404 for the response : %v", w.Code)
			}
			t.Log("\t [SUCCESS] Should receive a status code of 404 for the response.")
		}
	}
}

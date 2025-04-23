package elation

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllergyDocumentationService_Find(t *testing.T) {
	assert := assert.New(t)

	opts := &FindAllergiesDocumentationOptions{
		Pagination: &Pagination{
			Limit:  1,
			Offset: 2,
		},

		Patient: []int64{1},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/allergies_documentation", r.URL.Path)

		patient := r.URL.Query()["patient"]

		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")

		assert.Equal(opts.Patient, sliceStrToInt64(patient))
		assert.Equal(opts.Limit, strToInt(limit))
		assert.Equal(opts.Offset, strToInt(offset))

		b, err := json.Marshal(Response[[]*AllergyDocumentation]{
			Results: []*AllergyDocumentation{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			},
		})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := AllergyDocumentationService{client}

	found, res, err := svc.Find(context.Background(), opts)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestAllergyDocumentationService_Get(t *testing.T) {
	assert := assert.New(t)

	var id int64 = 1

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/allergies_documentation/"+strconv.FormatInt(id, 10), r.URL.Path)

		b, err := json.Marshal(&AllergyDocumentation{
			ID: id,
		})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := AllergyDocumentationService{client}

	found, res, err := svc.Get(context.Background(), id)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

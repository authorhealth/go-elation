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

func TestPhysicianService_Find(t *testing.T) {
	assert := assert.New(t)

	opts := &FindPhysiciansOptions{
		FirstName: "first name",
		LastName:  "last name",
		NPI:       "npi",
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/physicians", r.URL.Path)

		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")
		npi := r.URL.Query().Get("npi")

		assert.Equal(opts.FirstName, firstName)
		assert.Equal(opts.LastName, lastName)
		assert.Equal(opts.NPI, npi)

		b, err := json.Marshal(Response[[]*Physician]{
			Results: []*Physician{
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

	client := NewClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := PhysicianService{client}

	found, res, err := svc.Find(context.Background(), opts)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestPhysicianService_Get(t *testing.T) {
	assert := assert.New(t)

	var id int64 = 1

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/physicians/"+strconv.FormatInt(id, 10), r.URL.Path)

		b, err := json.Marshal(&Physician{
			ID: id,
		})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := PhysicianService{client}

	found, res, err := svc.Get(context.Background(), id)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

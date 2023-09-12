package elation

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceLocationService_Find(t *testing.T) {
	assert := assert.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/service_locations", r.URL.Path)

		b, err := json.Marshal(Response[[]*ServiceLocation]{
			Results: []*ServiceLocation{
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
		w.Write(b)
	}))
	defer srv.Close()

	client := NewClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := ServiceLocationService{client}

	found, res, err := svc.Find(context.Background())
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

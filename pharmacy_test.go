package elation

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPharmacyService_Get(t *testing.T) {
	assert := assert.New(t)

	ncpdpid := "1234789"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/pharmacies/"+ncpdpid, r.URL.Path)

		b, err := json.Marshal(&Pharmacy{
			ID:      123,
			NCPDPID: ncpdpid,
		})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := PharmacyService{client}

	found, res, err := svc.Get(context.Background(), ncpdpid)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

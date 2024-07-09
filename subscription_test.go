package elation

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSubscriptionService_Find(t *testing.T) {
	assert := assert.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/app/subscriptions/", r.URL.Path)

		b, err := json.Marshal(Response[[]*Subscription]{
			Results: []*Subscription{
				{
					ID:          1,
					CreatedDate: SubscriptionJSONDate(time.Now()),
				},
				{
					ID:          2,
					CreatedDate: SubscriptionJSONDate(time.Now()),
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
	svc := SubscriptionService{client}

	found, res, err := svc.Find(context.Background())
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestSubscriptionService_Subscribe(t *testing.T) {
	assert := assert.New(t)

	expected := &Subscribe{
		Resource:   "resource",
		Target:     "target",
		Properties: []byte(`{"key":"value"}`),
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodPost, r.Method)
		assert.Equal("/app/subscriptions/", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		assert.NoError(err)

		actual := &Subscribe{}
		err = json.Unmarshal(body, actual)
		assert.NoError(err)

		assert.Equal(expected.Resource, actual.Resource)
		assert.Equal(expected.Target, actual.Target)
		assert.Equal(expected.Properties, actual.Properties)

		b, err := json.Marshal(&Subscription{})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := SubscriptionService{client}

	found, res, err := svc.Subscribe(context.Background(), expected)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestSubscriptionService_Delete(t *testing.T) {
	assert := assert.New(t)

	var id int64 = 1

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodDelete, r.Method)
		assert.Equal("/app/subscriptions/"+strconv.FormatInt(id, 10)+"/", r.URL.Path)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := SubscriptionService{client}

	res, err := svc.Delete(context.Background(), id)
	assert.NotNil(res)
	assert.NoError(err)
}

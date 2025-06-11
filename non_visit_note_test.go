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

func TestNonVisitNoteService_Create(t *testing.T) {
	testCases := map[string]struct {
		create *NonVisitNoteCreate
	}{
		"required fields only request": {
			create: &NonVisitNoteCreate{
				Bullets: []*NonVisitNoteBullet{
					{
						Text:    "Patient",
						Version: 1,
						Author:  12345,
					},
				},
				ChartDate:    time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC),
				DocumentDate: time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC),
				Patient:      12345,
				Type:         "nonvisit",
			},
		},
		"signed fields request": {
			create: &NonVisitNoteCreate{
				Bullets: []*NonVisitNoteBullet{
					{
						Text:    "Patient",
						Version: 1,
						Author:  12345,
					},
				},
				ChartDate:    time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC),
				DocumentDate: time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC),
				Patient:      12345,
				Type:         "nonvisit",
				SignedBy:     12345,
				SignedDate:   Ptr(time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC)),
			},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tokenRequest(w, r) {
					return
				}

				assert.Equal(http.MethodPost, r.Method)
				assert.Equal("/non_visit_notes", r.URL.Path)

				body, err := io.ReadAll(r.Body)
				assert.NoError(err)

				create := &NonVisitNoteCreate{}
				err = json.Unmarshal(body, create)
				assert.NoError(err)

				assert.Equal(testCase.create, create)

				b, err := json.Marshal(&NonVisitNote{})
				assert.NoError(err)

				w.Header().Set("Content-Type", "application/json")
				//nolint
				w.Write(b)
			}))
			defer srv.Close()

			client := NewHTTPClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
			svc := NonVisitNoteService{client}

			created, res, err := svc.Create(context.Background(), testCase.create)
			assert.NotNil(created)
			assert.NotNil(res)
			assert.NoError(err)
		})
	}
}

func TestNonVisitNoteService_Find(t *testing.T) {
	assert := assert.New(t)

	opts := &FindNonVisitNotesOptions{
		Pagination: &Pagination{
			Limit:  1,
			Offset: 2,
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/non_visit_notes", r.URL.Path)

		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")

		assert.Equal(opts.Pagination.Limit, strToInt(limit))
		assert.Equal(opts.Pagination.Offset, strToInt(offset))

		b, err := json.Marshal(Response[[]*NonVisitNote]{
			Results: []*NonVisitNote{
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
	svc := NonVisitNoteService{client}

	found, res, err := svc.Find(context.Background(), opts)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestNonVisitNoteService_Get(t *testing.T) {
	assert := assert.New(t)

	var id int64 = 1

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/non_visit_notes/"+strconv.FormatInt(id, 10), r.URL.Path)

		b, err := json.Marshal(&NonVisitNote{
			ID: id,
		})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := NonVisitNoteService{client}

	found, res, err := svc.Get(context.Background(), id)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

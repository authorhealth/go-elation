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

func TestDiscontinuedMedicationService_Create(t *testing.T) {
	testCases := map[string]struct {
		create *DiscontinuedMedicationCreate
	}{
		"minimally-specified request": {
			create: &DiscontinuedMedicationCreate{
				MedOrder: 12345,
			},
		},
		"fully-specified request": {
			create: &DiscontinuedMedicationCreate{
				MedOrder:             12345,
				DiscontinueDate:      "2024-04-15",
				Reason:               "a very good reason",
				IsDocumented:         Ptr(true),
				DocumentingPersonnel: 67890,
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
				assert.Equal("/discontinued_medications", r.URL.Path)

				body, err := io.ReadAll(r.Body)
				assert.NoError(err)

				create := &DiscontinuedMedicationCreate{}
				err = json.Unmarshal(body, create)
				assert.NoError(err)

				assert.Equal(testCase.create, create)

				b, err := json.Marshal(&DiscontinuedMedication{})
				assert.NoError(err)

				w.Header().Set("Content-Type", "application/json")
				//nolint
				w.Write(b)
			}))
			defer srv.Close()

			client := NewHttpClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
			svc := DiscontinuedMedicationService{client}

			created, res, err := svc.Create(context.Background(), testCase.create)
			assert.NotNil(created)
			assert.NotNil(res)
			assert.NoError(err)
		})
	}
}

func TestDiscontinuedMedicationService_Find(t *testing.T) {
	assert := assert.New(t)

	opts := &FindDiscontinuedMedicationsOptions{
		Pagination: &Pagination{
			Limit:  1,
			Offset: 2,
		},

		Patient:            []int64{12345, 67890},
		Practice:           []int64{98765, 43210},
		DiscontinueDateLTE: time.Now(),
		DiscontinueDateGTE: time.Now(),
		DocumentDateLTE:    time.Now(),
		DocumentDateGTE:    time.Now(),
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/discontinued_medications", r.URL.Path)

		patient := r.URL.Query()["patient"]
		practice := r.URL.Query()["practice"]
		discontinueDateLTE := r.URL.Query().Get("discontinue_date__lte")
		discontinueDateGTE := r.URL.Query().Get("discontinue_date__gte")
		documentDateLTE := r.URL.Query().Get("document_date__lte")
		documentDateGTE := r.URL.Query().Get("document_date__gte")

		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")

		assert.Equal(opts.Patient, sliceStrToInt64(patient))
		assert.Equal(opts.Practice, sliceStrToInt64(practice))
		assert.Equal(opts.DiscontinueDateLTE.Format(time.RFC3339), discontinueDateLTE)
		assert.Equal(opts.DiscontinueDateGTE.Format(time.RFC3339), discontinueDateGTE)
		assert.Equal(opts.DocumentDateLTE.Format(time.RFC3339), documentDateLTE)
		assert.Equal(opts.DocumentDateGTE.Format(time.RFC3339), documentDateGTE)

		assert.Equal(opts.Pagination.Limit, strToInt(limit))
		assert.Equal(opts.Pagination.Offset, strToInt(offset))

		b, err := json.Marshal(Response[[]*DiscontinuedMedication]{
			Results: []*DiscontinuedMedication{
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

	client := NewHttpClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := DiscontinuedMedicationService{client}

	found, res, err := svc.Find(context.Background(), opts)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestDiscontinuedMedicationService_Get(t *testing.T) {
	assert := assert.New(t)

	var id int64 = 1

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/discontinued_medications/"+strconv.FormatInt(id, 10), r.URL.Path)

		b, err := json.Marshal(&DiscontinuedMedication{
			ID: id,
		})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewHttpClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := DiscontinuedMedicationService{client}

	found, res, err := svc.Get(context.Background(), id)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestDiscontinuedMedicationService_Update(t *testing.T) {
	assert := assert.New(t)

	var id int64 = 1
	expected := &DiscontinuedMedicationUpdate{
		DiscontinueDate:      "2024-04-15",
		Reason:               "a very good reason",
		DocumentingPersonnel: 67890,
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodPatch, r.Method)
		assert.Equal("/discontinued_medications/"+strconv.FormatInt(id, 10), r.URL.Path)

		body, err := io.ReadAll(r.Body)
		assert.NoError(err)

		actual := &DiscontinuedMedicationUpdate{}
		err = json.Unmarshal(body, actual)
		assert.NoError(err)

		assert.Equal(expected, actual)

		b, err := json.Marshal(&DiscontinuedMedication{})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewHttpClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := DiscontinuedMedicationService{client}

	updated, res, err := svc.Update(context.Background(), id, expected)
	assert.NotNil(updated)
	assert.NotNil(res)
	assert.NoError(err)
}

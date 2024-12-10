package elation

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVisitNoteService_Create(t *testing.T) {
	testCases := map[string]struct {
		create *VisitNoteCreate
	}{
		"required fields only request": {
			create: &VisitNoteCreate{
				Bullets: []*VisitNoteBullet{
					{
						Category: "Subjective",
						Text:     "Patient",
						Version:  1,
						Sequence: 1,
						Author:   12345,
					},
				},
				ChartDate:    time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC),
				DocumentDate: time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC),
				Patient:      12345,
				Template:     "SOAP",
				Physician:    12345,
			},
		},
		"all specified fields request": {
			create: &VisitNoteCreate{
				Bullets: []*VisitNoteBullet{
					{
						Category: "Subjective",
						Text:     "Patient",
						Version:  1,
						Sequence: 1,
						Author:   12345,
					},
				},
				ChartDate:    time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC),
				DocumentDate: time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC),
				Patient:      12345,
				Template:     "SOAP",
				Physician:    12345,
				Type:         "Office Visit Note",
				Confidential: true,
				SignedBy:     12345,
				Signatures: []*VisitNoteSignature{
					{
						User:       12345,
						UserName:   "John Doe",
						SignedDate: time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC),
						Role:       "Physician",
						Comments:   new(string),
					},
				},
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
				assert.Equal("/visit_notes", r.URL.Path)

				body, err := io.ReadAll(r.Body)
				assert.NoError(err)

				create := &VisitNoteCreate{}
				err = json.Unmarshal(body, create)
				assert.NoError(err)

				assert.Equal(testCase.create, create)

				b, err := json.Marshal(&VisitNote{})
				assert.NoError(err)

				w.Header().Set("Content-Type", "application/json")
				//nolint
				w.Write(b)
			}))
			defer srv.Close()

			client := NewHTTPClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
			svc := VisitNoteService{client}

			created, res, err := svc.Create(context.Background(), testCase.create)
			assert.NotNil(created)
			assert.NotNil(res)
			assert.NoError(err)
		})
	}
}

func TestVisitNoteService_Find(t *testing.T) {
	assert := assert.New(t)

	opts := &FindVisitNotesOptions{
		Patient:         123,
		Physician:       456,
		Practice:        789,
		LastModifiedGT:  time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		LastModifiedGTE: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
		LastModifiedLT:  time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC),
		LastModifiedLTE: time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC),
		FromSignedDate:  time.Date(2022, 1, 5, 0, 0, 0, 0, time.UTC),
		ToSignedDate:    time.Date(2022, 1, 6, 0, 0, 0, 0, time.UTC),
		Unsigned:        true,
	}

	visitNotes := []*VisitNote{
		{
			ID: 1,
		},
		{
			ID: 2,
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/visit_notes", r.URL.Path)

		actualPatient := r.URL.Query().Get("patient")
		assert.Equal(opts.Patient, strToInt64(actualPatient))

		actualPhysician := r.URL.Query().Get("physician")
		assert.Equal(opts.Physician, strToInt64(actualPhysician))

		actualPractice := r.URL.Query().Get("practice")
		assert.Equal(opts.Practice, strToInt64(actualPractice))

		actualLastModifiedGT := r.URL.Query().Get("last_modified_gt")
		assert.Equal(opts.LastModifiedGT.Format(time.RFC3339), actualLastModifiedGT)

		actualLastModifiedGTE := r.URL.Query().Get("last_modified_gte")
		assert.Equal(opts.LastModifiedGTE.Format(time.RFC3339), actualLastModifiedGTE)

		actualLastModifiedLT := r.URL.Query().Get("last_modified_lt")
		assert.Equal(opts.LastModifiedLT.Format(time.RFC3339), actualLastModifiedLT)

		actualLastModifiedLTE := r.URL.Query().Get("last_modified_lte")
		assert.Equal(opts.LastModifiedLTE.Format(time.RFC3339), actualLastModifiedLTE)

		actualFromSignedDate := r.URL.Query().Get("from_signed_date")
		assert.Equal(opts.FromSignedDate.Format(time.RFC3339), actualFromSignedDate)

		actualToSignedDate := r.URL.Query().Get("to_signed_date")
		assert.Equal(opts.ToSignedDate.Format(time.RFC3339), actualToSignedDate)

		actualUnsigned := r.URL.Query().Get("unsigned")
		assert.Equal(opts.Unsigned, strToBool(actualUnsigned))

		b, err := json.Marshal(Response[[]*VisitNote]{
			Results: []*VisitNote{
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
	svc := VisitNoteService{client}

	visitNotesRes, res, err := svc.Find(context.Background(), opts)
	assert.NotEmpty(visitNotesRes)
	assert.NotNil(res)
	assert.NoError(err)

	assert.Equal(visitNotes, visitNotesRes.Results)
}

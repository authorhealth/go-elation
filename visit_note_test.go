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
						Category:    "Subjective",
						Text:        "Patient",
						Version:     1,
						Sequence:    1,
						Author:      12345,
						UpdatedDate: time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC),
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
						Category:    "Subjective",
						Text:        "Patient",
						Version:     1,
						Sequence:    1,
						Author:      12345,
						UpdatedDate: time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC),
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
				SignedDate:   time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC),
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

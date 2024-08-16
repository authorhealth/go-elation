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

func TestBillService_Create(t *testing.T) {
	testCases := map[string]struct {
		create *BillCreate
	}{
		"required fields only request": {
			create: &BillCreate{
				ServiceLocation: &BillServiceLocation{
					ID:                 1,
					Name:               "service location",
					IsPrimary:          true,
					PlaceOfService:     "office",
					PlaceOfServiceCode: "11",
					AddressLine1:       "123 Main St",
					City:               "Schenectady",
					State:              "NY",
					Zip:                "12345",
					Phone:              "555-234-5678",
					CreatedDate:        time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC),
				},
				VisitNoteID: 64409108504,
				Patient:     64901939201,
				Practice:    65540,
				Physician:   64811630594,
			},
		},
		"all specified fields request": {
			create: &BillCreate{
				ServiceLocation: &BillServiceLocation{
					ID:                 1,
					Name:               "service location",
					IsPrimary:          true,
					PlaceOfService:     "office",
					PlaceOfServiceCode: "11",
					AddressLine1:       "123 Main St",
					City:               "Schenectady",
					State:              "NY",
					Zip:                "12345",
					Phone:              "555-234-5678",
					CreatedDate:        time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC),
				},
				VisitNoteID: 64409108504,
				Patient:     64901939201,
				Practice:    65540,
				Physician:   64811630594,
				CPTs: []*BillCPT{
					{
						CPT:       12,
						Modifiers: []string{"modifier 1", "modifier 2"},
						DXs:       []string{"dx 1", "dx 2"},
						AltDXs:    []string{"alt dx 1", "alt dx 2"},
					},
				},
				BillingProvider:     42120898,
				RenderingProvider:   68382673,
				SupervisingProvider: 52893234,
				ReferringProvider: &BillProvider{
					Name:  "referring provider",
					State: "NY",
					NPI:   "1234567890",
				},
				OrderingProvider: &BillProvider{
					Name:  "referring provider",
					State: "NY",
					NPI:   "1234567890",
				},
				PriorAuthorization: "1234-ABC",
				PaymentAmount:      100.00,
				Notes:              "additional billing notes",
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
				assert.Equal("/bills", r.URL.Path)

				body, err := io.ReadAll(r.Body)
				assert.NoError(err)

				create := &BillCreate{}
				err = json.Unmarshal(body, create)
				assert.NoError(err)

				assert.Equal(testCase.create, create)

				b, err := json.Marshal(&Bill{})
				assert.NoError(err)

				w.Header().Set("Content-Type", "application/json")
				//nolint
				w.Write(b)
			}))
			defer srv.Close()

			client := NewHTTPClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
			svc := BillService{client}

			created, res, err := svc.Create(context.Background(), testCase.create)
			assert.NotNil(created)
			assert.NotNil(res)
			assert.NoError(err)
		})
	}
}

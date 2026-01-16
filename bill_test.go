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

func TestBillService_Create(t *testing.T) {
	testCases := map[string]struct {
		create *BillCreate
	}{
		"required fields only request": {
			create: &BillCreate{
				ServiceLocation: 10,
				VisitNote:       64409108504,
				Patient:         64901939201,
				Practice:        65540,
				Physician:       64811630594,
				CPTs:            []*CreatedBillCPT{},
			},
		},
		"all specified fields request": {
			create: &BillCreate{
				ServiceLocation: 10,
				VisitNote:       64409108504,
				Patient:         64901939201,
				Practice:        65540,
				Physician:       64811630594,
				CPTs: []*CreatedBillCPT{
					{
						CPT:        "12",
						Units:      "1.0",
						UnitCharge: "122.0",
						Modifiers:  []string{"modifier 1", "modifier 2"},
						DXs:        []CreatedBillDX{{ICD10Code: "dx 1"}, {ICD10Code: "dx 2"}},
						AltDXs:     []string{"alt dx 1", "alt dx 2"},
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

				b, err := json.Marshal(&CreatedBill{})
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

func TestBillService_Create_already_exists(t *testing.T) {
	assert := assert.New(t)

	billCreate := &BillCreate{
		ServiceLocation: 10,
		VisitNote:       64409108504,
		Patient:         64901939201,
		Practice:        65540,
		Physician:       64811630594,
		CPTs:            []*CreatedBillCPT{},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodPost, r.Method)
		assert.Equal("/bills", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		assert.NoError(err)

		actualBillCreate := &BillCreate{}
		err = json.Unmarshal(body, actualBillCreate)
		assert.NoError(err)

		assert.Equal(billCreate, actualBillCreate)

		errorRes := map[string][]string{
			"visit_note": {billExistError},
		}
		b, err := json.Marshal(errorRes)
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := BillService{client}

	created, res, err := svc.Create(context.Background(), billCreate)
	assert.Nil(created)
	assert.NotNil(res)
	assert.ErrorIs(err, ErrBillExist)
}

func TestBillService_Get(t *testing.T) {
	assert := assert.New(t)

	var id int64 = 65099661468

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/bills/"+strconv.FormatInt(id, 10), r.URL.Path)

		b, err := json.Marshal(&Bill{
			ID: id,
		})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := BillService{client}

	found, res, err := svc.Get(context.Background(), id)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
	assert.Equal(id, found.ID)
}

func TestBillService_Find(t *testing.T) {
	assert := assert.New(t)

	fromServiceDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	toServiceDate := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)

	opts := &FindBillOptions{
		Patient:           []int64{141289139994622},
		BillID:            []int64{143945720594544},
		FromServiceDate:   fromServiceDate,
		ToServiceDate:     toServiceDate,
		AssignedPhysician: []int64{64811630594},
		SigningPhysician:  64811630594,
		VisitNoteID:       []int64{12314823928},
	}

	bills := []*Bill{
		{ID: 101},
		{ID: 102},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/bills", r.URL.Path)

		q := r.URL.Query()
		assert.Equal("141289139994622", q.Get("patient"))
		assert.Equal("143945720594544", q.Get("bill_id"))
		assert.Equal(fromServiceDate.Format(time.RFC3339), q.Get("from_service_date"))
		assert.Equal(toServiceDate.Format(time.RFC3339), q.Get("to_service_date"))
		assert.Equal("64811630594", q.Get("assigned_physician"))
		assert.Equal("64811630594", q.Get("signing_physician"))
		assert.Equal("12314823928", q.Get("visit_note_id"))

		b, err := json.Marshal(Response[[]*Bill]{
			Results: bills,
		})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := BillService{client}

	billsRes, res, err := svc.Find(context.Background(), opts)
	assert.NotNil(billsRes)
	assert.NotNil(res)
	assert.NoError(err)

	assert.Len(billsRes.Results, 2)
	assert.Equal(bills, billsRes.Results)
}

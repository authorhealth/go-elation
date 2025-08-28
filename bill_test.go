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
				CPTs:            []*BillCPTCreate{},
			},
		},
		"all specified fields request": {
			create: &BillCreate{
				ServiceLocation: 10,
				VisitNote:       64409108504,
				Patient:         64901939201,
				Practice:        65540,
				Physician:       64811630594,
				CPTs: []*BillCPTCreate{
					{
						CPT:        "12",
						Units:      "1.0",
						UnitCharge: "122.0",
						Modifiers:  []string{"modifier 1", "modifier 2"},
						DXs:        []BillDX{{ICD10Code: "dx 1"}, {ICD10Code: "dx 2"}},
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

func TestBillService_Create_already_exists(t *testing.T) {
	assert := assert.New(t)

	billCreate := &BillCreate{
		ServiceLocation: 10,
		VisitNote:       64409108504,
		Patient:         64901939201,
		Practice:        65540,
		Physician:       64811630594,
		CPTs:            []*BillCPTCreate{},
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

func TestBillService_Find(t *testing.T) {
	assert := assert.New(t)

	opts := &FindBillsOptions{
		Pagination: &Pagination{
			Limit:  1,
			Offset: 2,
		},

		AssignedPhysician: []int64{1},
		BillID:            []int64{2},
		FromServiceDate:   time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
		ToServiceDate:     time.Date(2023, 5, 20, 0, 0, 0, 0, time.UTC),
		Patient:           []int64{3},
		Practice:          []int64{4},
		SigningPhysician:  []int64{5},
		VisitNoteID:       []int64{6},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/bills", r.URL.Path)

		assignedPhysician := r.URL.Query().Get("assigned_physician")
		billID := r.URL.Query().Get("bill_id")
		fromServiceDate := r.URL.Query().Get("from_service_date")
		toServiceDate := r.URL.Query().Get("to_service_date")
		patient := r.URL.Query().Get("patient")
		practice := r.URL.Query().Get("practice")
		signingPhysician := r.URL.Query().Get("signing_physician")
		visitNoteID := r.URL.Query().Get("visit_note_id")

		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")

		assert.Equal(opts.AssignedPhysician, sliceStrToInt64([]string{assignedPhysician}))
		assert.Equal(opts.BillID, sliceStrToInt64([]string{billID}))
		assert.Equal(opts.FromServiceDate.Format(time.RFC3339), fromServiceDate)
		assert.Equal(opts.ToServiceDate.Format(time.RFC3339), toServiceDate)
		assert.Equal(opts.Patient, sliceStrToInt64([]string{patient}))
		assert.Equal(opts.Practice, sliceStrToInt64([]string{practice}))
		assert.Equal(opts.SigningPhysician, sliceStrToInt64([]string{signingPhysician}))
		assert.Equal(opts.VisitNoteID, sliceStrToInt64([]string{visitNoteID}))

		assert.Equal(opts.Limit, strToInt(limit))
		assert.Equal(opts.Offset, strToInt(offset))

		b, err := json.Marshal(Response[[]*Bill]{
			Results: []*Bill{
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
	svc := BillService{client}

	found, res, err := svc.Find(context.Background(), opts)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestBillService_Get(t *testing.T) {
	assert := assert.New(t)

	var id int64 = 1

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/bills/"+strconv.FormatInt(id, 10), r.URL.Path)

		b, err := json.Marshal(&Patient{
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
}

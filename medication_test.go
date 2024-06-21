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

func TestMedicationService_Create(t *testing.T) {
	testCases := map[string]struct {
		create *PatientMedicationCreate
	}{
		"minimally-specified request": {
			create: &PatientMedicationCreate{},
		},
		"fully-specified request": {
			create: &PatientMedicationCreate{
				AuthRefills: 5,
				IsDocMed:    true,
				Medication: &PatientMedicationCreateMedication{
					ID:         12345,
					NDCs:       []string{"ndc 1", "ndc 2"},
					RxnormCuis: []string{"rxnormCui 1", "rxnormCui 2"},
				},
				MedicationType:       "otc",
				OrderType:            "New",
				Patient:              12345,
				Practice:             12345,
				PrescribingPhysician: 12345,
				Qty:                  "20",
				QtyUnits:             "tab",
				StartDate:            "2024-04-15",
				Icd10Codes: []*PatientMedicationCreateICD10Code{
					{Code: "icd-10 code 1"},
					{Code: "icd-10 code 2"},
				},
				Thread: &PatientMedicationCreateThread{
					ID:          12345,
					IsPermanent: true,
				},
				Directions:           "take 1 tab with water twice a day",
				Notes:                "some notes",
				PharmacyInstructions: "some instructions to the pharmacy",
				NumSamples:           "5",
				DocumentDate:         Ptr(time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC)),
				DocumentingPersonnel: 12345,
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
				assert.Equal("/medications", r.URL.Path)

				body, err := io.ReadAll(r.Body)
				assert.NoError(err)

				create := &PatientMedicationCreate{}
				err = json.Unmarshal(body, create)
				assert.NoError(err)

				assert.Equal(testCase.create, create)

				b, err := json.Marshal(&Medication{})
				assert.NoError(err)

				w.Header().Set("Content-Type", "application/json")
				//nolint
				w.Write(b)
			}))
			defer srv.Close()

			client := NewHttpClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
			svc := MedicationService{client}

			created, res, err := svc.Create(context.Background(), testCase.create)
			assert.NotNil(created)
			assert.NotNil(res)
			assert.NoError(err)
		})
	}
}

func TestMedicationService_Find(t *testing.T) {
	assert := assert.New(t)

	opts := &FindPatientMedicationsOptions{
		Pagination: &Pagination{
			Limit:  1,
			Offset: 2,
		},

		Patient:         1,
		Thread:          2,
		Practice:        3,
		Discontinued:    true,
		Permanent:       true,
		DocumentDateLT:  time.Now(),
		DocumentDateGT:  time.Now(),
		DocumentDateLTE: time.Now(),
		DocumentDateGTE: time.Now(),
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/medications", r.URL.Path)

		patient := r.URL.Query().Get("patient")
		thread := r.URL.Query().Get("thread")
		practice := r.URL.Query().Get("practice")
		discontinued := r.URL.Query().Get("discontinued")
		permanent := r.URL.Query().Get("permanent")
		documentDateLT := r.URL.Query().Get("document_date__lt")
		documentDateGT := r.URL.Query().Get("document_date__gt")
		documentDateLTE := r.URL.Query().Get("document_date__lte")
		documentDateGTE := r.URL.Query().Get("document_date__gte")

		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")

		assert.Equal(opts.Patient, strToInt64(patient))
		assert.Equal(opts.Thread, strToInt64(thread))
		assert.Equal(opts.Practice, strToInt64(practice))
		assert.Equal(opts.Discontinued, strToBool(discontinued))
		assert.Equal(opts.Permanent, strToBool(permanent))
		assert.Equal(opts.DocumentDateLT.Format(time.RFC3339), documentDateLT)
		assert.Equal(opts.DocumentDateGT.Format(time.RFC3339), documentDateGT)
		assert.Equal(opts.DocumentDateLTE.Format(time.RFC3339), documentDateLTE)
		assert.Equal(opts.DocumentDateGTE.Format(time.RFC3339), documentDateGTE)

		assert.Equal(opts.Pagination.Limit, strToInt(limit))
		assert.Equal(opts.Pagination.Offset, strToInt(offset))

		b, err := json.Marshal(Response[[]*PatientMedication]{
			Results: []*PatientMedication{
				{
					ID:      1,
					Patient: 2,
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
	svc := MedicationService{client}

	found, res, err := svc.Find(context.Background(), opts)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestMedicationService_Get(t *testing.T) {
	assert := assert.New(t)

	var id int64 = 1

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/medications/"+strconv.FormatInt(id, 10), r.URL.Path)

		b, err := json.Marshal(&Physician{
			ID: id,
		})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewHttpClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := MedicationService{client}

	found, res, err := svc.Get(context.Background(), id)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

package elation

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsuranceEligibilityService_Create(t *testing.T) {
	testCases := map[string]struct {
		create *InsuranceEligibilityCreate
	}{
		"minimally-specified request": {
			create: &InsuranceEligibilityCreate{},
		},
		"fully-specified request": {
			create: &InsuranceEligibilityCreate{
				GroupID:   12345,
				GroupName: "group name",
				GroupNPI:  "group NPI",
			},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			var id int64 = 1

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tokenRequest(w, r) {
					return
				}

				assert.Equal(http.MethodPost, r.Method)
				assert.Equal("/patient_insurances/1/eligibility/", r.URL.Path)

				body, err := io.ReadAll(r.Body)
				assert.NoError(err)

				create := &InsuranceEligibilityCreate{}
				err = json.Unmarshal(body, create)
				assert.NoError(err)

				assert.Equal(testCase.create, create)

				b, err := json.Marshal(&InsuranceEligibility{})
				assert.NoError(err)

				w.Header().Set("Content-Type", "application/json")
				//nolint
				w.Write(b)
			}))
			defer srv.Close()

			client := NewHTTPClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
			svc := InsuranceEligibilityService{client}

			created, res, err := svc.Create(context.Background(), id, testCase.create)
			assert.NotNil(created)
			assert.NotNil(res)
			assert.NoError(err)
		})
	}
}

func TestInsuranceEligibilityService_Get(t *testing.T) {
	assert := assert.New(t)

	var id int64 = 1

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/patient_insurances/1/eligibility/", r.URL.Path)

		b, err := json.Marshal(&InsuranceEligibility{})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := InsuranceEligibilityService{client}

	found, res, err := svc.Get(context.Background(), id)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestInsuranceEligibilityService_GetFullReport(t *testing.T) {
	assert := assert.New(t)

	var id int64 = 1

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/patient_insurances/1/eligibility_full_report/", r.URL.Path)

		b, err := json.Marshal(&InsuranceEligibility{})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := InsuranceEligibilityService{client}

	found, res, err := svc.GetFullReport(context.Background(), id)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

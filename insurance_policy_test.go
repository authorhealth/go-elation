package elation

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsurancePolicyService_Create(t *testing.T) {
	assert := assert.New(t)

	var patientID int64 = 500

	expectedCreate := &InsurancePolicyCreate{
		PracticeID:            100,
		PatientID:             patientID,
		PatientFirstName:      "Alice",
		PatientLastName:       "Smith",
		PatientDOB:            "1955-01-01",
		Status:                "active",
		Rank:                  Ptr[int64](1),
		MemberID:              Ptr("MEMBER123"),
		RelationshipToInsured: Ptr("self"),
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodPost, r.Method)
		assert.Equal("/patients/"+strconv.FormatInt(patientID, 10)+"/policies", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		assert.NoError(err)

		create := &InsurancePolicyCreate{}
		err = json.Unmarshal(body, create)
		assert.NoError(err)

		assert.Equal(expectedCreate.PatientFirstName, create.PatientFirstName)
		assert.Equal(expectedCreate.PatientLastName, create.PatientLastName)
		assert.Equal(expectedCreate.PatientDOB, create.PatientDOB)
		assert.Equal(*expectedCreate.MemberID, *create.MemberID)
		assert.Equal(*expectedCreate.Rank, *create.Rank)
		assert.Equal(expectedCreate.Status, create.Status)
		assert.Equal(expectedCreate.PatientID, create.PatientID)
		assert.Equal(expectedCreate.PracticeID, create.PracticeID)

		b, err := json.Marshal(&InsurancePolicy{ID: 1})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := InsurancePolicyService{client}

	created, res, err := svc.Create(context.Background(), patientID, expectedCreate)
	assert.NotNil(created)
	assert.NotNil(res)
	assert.NoError(err)
	assert.Equal(int64(1), created.ID)
}

func TestInsurancePolicyService_Find(t *testing.T) {
	assert := assert.New(t)

	var patientID int64 = 500

	opts := &FindInsurancePoliciesOptions{
		ActiveOnly: Ptr(true),
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/patients/"+strconv.FormatInt(patientID, 10)+"/policies", r.URL.Path)

		// Verify query parameters
		activeOnlyStr := r.URL.Query().Get("active_only")
		activeOnly, err := strconv.ParseBool(activeOnlyStr)
		assert.NoError(err)

		assert.Equal(*opts.ActiveOnly, activeOnly)

		b, err := json.Marshal(FindInsurancePoliciesResponse{
			Results: []*InsurancePolicy{
				{ID: 101, Status: "active", PatientID: patientID, Rank: Ptr[int64](1)},
				{ID: 102, Status: "inactive", PatientID: patientID, Rank: Ptr[int64](4)},
			},
		})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := InsurancePolicyService{client}

	found, res, err := svc.Find(context.Background(), patientID, opts)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
	assert.Len(found.Results, 2)
	assert.Equal(int64(101), found.Results[0].ID)
}

func TestInsurancePolicyService_Get(t *testing.T) {
	assert := assert.New(t)

	var patientID int64 = 500
	var id int64 = 99

	expectedPolicy := &InsurancePolicy{
		ID:        id,
		Status:    "inactive",
		PatientID: patientID,
		MemberID:  Ptr("MEMBER123"),
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/patients/"+strconv.FormatInt(patientID, 10)+"/policies/"+strconv.FormatInt(id, 10), r.URL.Path)

		b, err := json.Marshal(expectedPolicy)
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := InsurancePolicyService{client}

	found, res, err := svc.Get(context.Background(), patientID, id)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
	assert.Equal(id, found.ID)
	assert.Equal("inactive", found.Status)
	assert.Equal("MEMBER123", *found.MemberID)
}

func TestInsurancePolicyService_Update(t *testing.T) {
	assert := assert.New(t)

	var patientID int64 = 500
	var id int64 = 42

	expectedUpdate := &InsurancePolicyUpdate{
		ID:     id,
		Status: "active",
		Copay:  Ptr("50.00"),
		Rank:   Ptr[int64](2),
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodPut, r.Method)
		assert.Equal("/patients/"+strconv.FormatInt(patientID, 10)+"/policies/"+strconv.FormatInt(id, 10), r.URL.Path)

		body, err := io.ReadAll(r.Body)
		assert.NoError(err)

		actual := &InsurancePolicyUpdate{}
		err = json.Unmarshal(body, actual)
		assert.NoError(err)

		assert.Equal(expectedUpdate.Status, actual.Status)
		assert.Equal(*expectedUpdate.Copay, *actual.Copay)
		assert.Equal(*expectedUpdate.Rank, *actual.Rank)

		b, err := json.Marshal(&InsurancePolicy{ID: id})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := InsurancePolicyService{client}

	updated, res, err := svc.Update(context.Background(), patientID, id, expectedUpdate)
	assert.NotNil(updated)
	assert.NotNil(res)
	assert.NoError(err)
	assert.Equal(id, updated.ID)
}

func TestInsurancePolicyService_Delete(t *testing.T) {
	assert := assert.New(t)

	var patientID int64 = 500
	var id int64 = 77

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodDelete, r.Method)
		assert.Equal("/patients/"+strconv.FormatInt(patientID, 10)+"/policies/"+strconv.FormatInt(id, 10), r.URL.Path)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := InsurancePolicyService{client}

	res, err := svc.Delete(context.Background(), patientID, id)
	assert.NotNil(res)
	assert.NoError(err)
	assert.Equal(http.StatusNoContent, res.StatusCode)
}

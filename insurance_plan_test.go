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

func TestInsurancePlanService_Create(t *testing.T) {
	assert := assert.New(t)

	expectedCreate := &InsurancePlanCreate{
		Practice:         12345,
		InsuranceCompany: 67890,
		Name:             "plan name",
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodPost, r.Method)
		assert.Equal("/insurance_plans", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		assert.NoError(err)

		create := &InsurancePlanCreate{}
		err = json.Unmarshal(body, create)
		assert.NoError(err)

		assert.Equal(expectedCreate, create)

		b, err := json.Marshal(&InsurancePlan{})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewHttpClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := InsurancePlanService{client}

	created, res, err := svc.Create(context.Background(), expectedCreate)
	assert.NotNil(created)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestInsurancePlanService_Find(t *testing.T) {
	assert := assert.New(t)

	opts := &FindInsurancePlansOptions{
		Pagination: &Pagination{
			Limit:  1,
			Offset: 2,
		},

		Practice:         []int64{12345, 67890},
		InsuranceCompany: []int64{98765, 43210},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/insurance_plans", r.URL.Path)

		practice := r.URL.Query()["practice"]
		insuranceCompany := r.URL.Query()["insurance_company"]

		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")

		assert.Equal(opts.Practice, sliceStrToInt64(practice))
		assert.Equal(opts.InsuranceCompany, sliceStrToInt64(insuranceCompany))

		assert.Equal(opts.Pagination.Limit, strToInt(limit))
		assert.Equal(opts.Pagination.Offset, strToInt(offset))

		b, err := json.Marshal(Response[[]*InsurancePlan]{
			Results: []*InsurancePlan{
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
	svc := InsurancePlanService{client}

	found, res, err := svc.Find(context.Background(), opts)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestInsurancePlanService_Get(t *testing.T) {
	assert := assert.New(t)

	var id int64 = 1

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/insurance_plans/"+strconv.FormatInt(id, 10), r.URL.Path)

		b, err := json.Marshal(&InsurancePlan{
			ID: id,
		})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewHttpClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := InsurancePlanService{client}

	found, res, err := svc.Get(context.Background(), id)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestInsurancePlanService_Update(t *testing.T) {
	assert := assert.New(t)

	var id int64 = 1
	expected := &InsurancePlanUpdate{
		Name: "plan name",
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodPatch, r.Method)
		assert.Equal("/insurance_plans/"+strconv.FormatInt(id, 10), r.URL.Path)

		body, err := io.ReadAll(r.Body)
		assert.NoError(err)

		actual := &InsurancePlanUpdate{}
		err = json.Unmarshal(body, actual)
		assert.NoError(err)

		assert.Equal(expected, actual)

		b, err := json.Marshal(&InsurancePlan{})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewHttpClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := InsurancePlanService{client}

	updated, res, err := svc.Update(context.Background(), id, expected)
	assert.NotNil(updated)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestInsurancePlanService_Delete(t *testing.T) {
	assert := assert.New(t)

	var id int64 = 1

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodDelete, r.Method)
		assert.Equal("/insurance_plans/"+strconv.FormatInt(id, 10), r.URL.Path)
	}))
	defer srv.Close()

	client := NewHttpClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := InsurancePlanService{client}

	res, err := svc.Delete(context.Background(), id)
	assert.NotNil(res)
	assert.NoError(err)
}

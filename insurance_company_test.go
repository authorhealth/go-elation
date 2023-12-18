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

func TestInsuranceCompanyService_Create(t *testing.T) {
	testCases := map[string]struct {
		create *InsuranceCompanyCreate
	}{
		"minimally-specified request": {
			create: &InsuranceCompanyCreate{
				Practice: 12345,
				Carrier:  "Insurance Co.",
			},
		},
		"fully-specified request": {
			create: &InsuranceCompanyCreate{
				Practice:         12345,
				Carrier:          "Insurance Co.",
				Address:          "123 Any St",
				Suite:            "Unit 5B",
				City:             "Schenectady",
				State:            "NY",
				Zip:              "12345",
				Phone:            "555-234-5678",
				Extension:        "x1234",
				PayerID:          "67890",
				ExternalVendorID: "ABC123",
				IsCredentialed:   Ptr(true),
				Aliases:          "InsuranceCo., The Grand Insurance Company",
				InsuranceType:    "commercial",
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
				assert.Equal("/insurance_companies", r.URL.Path)

				body, err := io.ReadAll(r.Body)
				assert.NoError(err)

				create := &InsuranceCompanyCreate{}
				err = json.Unmarshal(body, create)
				assert.NoError(err)

				assert.Equal(testCase.create, create)

				b, err := json.Marshal(&InsuranceCompany{})
				assert.NoError(err)

				w.Header().Set("Content-Type", "application/json")
				//nolint
				w.Write(b)
			}))
			defer srv.Close()

			client := NewClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
			svc := InsuranceCompanyService{client}

			created, res, err := svc.Create(context.Background(), testCase.create)
			assert.NotNil(created)
			assert.NotNil(res)
			assert.NoError(err)
		})
	}
}

func TestInsuranceCompanyService_Find(t *testing.T) {
	assert := assert.New(t)

	opts := &FindInsuranceCompaniesOptions{
		Pagination: &Pagination{
			Limit:  1,
			Offset: 2,
		},

		Practice: []int64{12345, 67890},
		Carrier:  "Insurance Co.",
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/insurance_companies", r.URL.Path)

		practice := r.URL.Query()["practice"]
		carrier := r.URL.Query().Get("carrier")

		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")

		assert.Equal(opts.Practice, sliceStrToInt64(practice))
		assert.Equal(opts.Carrier, carrier)

		assert.Equal(opts.Pagination.Limit, strToInt(limit))
		assert.Equal(opts.Pagination.Offset, strToInt(offset))

		b, err := json.Marshal(Response[[]*InsuranceCompany]{
			Results: []*InsuranceCompany{
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

	client := NewClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := InsuranceCompanyService{client}

	found, res, err := svc.Find(context.Background(), opts)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestInsuranceCompanyService_Get(t *testing.T) {
	assert := assert.New(t)

	var id int64 = 1

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/insurance_companies/"+strconv.FormatInt(id, 10), r.URL.Path)

		b, err := json.Marshal(&InsuranceCompany{
			ID: id,
		})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := InsuranceCompanyService{client}

	found, res, err := svc.Get(context.Background(), id)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestInsuranceCompanyService_Update(t *testing.T) {
	assert := assert.New(t)

	var id int64 = 1
	expected := &InsuranceCompanyUpdate{
		Carrier:          "Insurance Co.",
		Address:          "123 Any St",
		Suite:            "Unit 5B",
		City:             "Schenectady",
		State:            "NY",
		Zip:              "12345",
		Phone:            "555-234-5678",
		Extension:        "x1234",
		PayerID:          "67890",
		ExternalVendorID: "ABC123",
		IsCredentialed:   Ptr(true),
		Aliases:          "InsuranceCo., The Grand Insurance Company",
		InsuranceType:    "commercial",
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodPatch, r.Method)
		assert.Equal("/insurance_companies/"+strconv.FormatInt(id, 10), r.URL.Path)

		body, err := io.ReadAll(r.Body)
		assert.NoError(err)

		actual := &InsuranceCompanyUpdate{}
		err = json.Unmarshal(body, actual)
		assert.NoError(err)

		assert.Equal(expected, actual)

		b, err := json.Marshal(&InsuranceCompany{})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := InsuranceCompanyService{client}

	updated, res, err := svc.Update(context.Background(), id, expected)
	assert.NotNil(updated)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestInsuranceCompanyService_Delete(t *testing.T) {
	assert := assert.New(t)

	var id int64 = 1

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodDelete, r.Method)
		assert.Equal("/insurance_companies/"+strconv.FormatInt(id, 10), r.URL.Path)
	}))
	defer srv.Close()

	client := NewClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := InsuranceCompanyService{client}

	res, err := svc.Delete(context.Background(), id)
	assert.NotNil(res)
	assert.NoError(err)
}

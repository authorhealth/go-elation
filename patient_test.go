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

	"cloud.google.com/go/civil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func TestPatientService_Create(t *testing.T) {
	testCases := map[string]struct {
		create *PatientCreate
	}{
		"minimally-specified request": {
			create: &PatientCreate{
				LastName:          "last name",
				FirstName:         "first name",
				Sex:               "sex",
				DOB:               "dob",
				PrimaryPhysician:  1,
				CaregiverPractice: 2,
			},
		},
		"fully-specified request": {
			create: &PatientCreate{
				LastName:          "last name",
				FirstName:         "first name",
				Sex:               "sex",
				DOB:               "dob",
				PrimaryPhysician:  1,
				CaregiverPractice: 2,
				Address: &PatientAddress{
					AddressLine1: "123 Any St",
					AddressLine2: "Unit 5B",
					City:         "Schenectady",
					State:        "NY",
					Zip:          "12345",
				},
				Phones: []*PatientPhone{
					{
						Phone:     "555-234-5678",
						PhoneType: "Mobile",
					},
					{
						Phone:     "555-987-6543",
						PhoneType: "Home",
					},
				},
				Emails: []*PatientEmail{
					{
						Email: "x@y.net",
					},
					{
						Email: "a@b.com",
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
				assert.Equal("/patients", r.URL.Path)

				body, err := io.ReadAll(r.Body)
				assert.NoError(err)

				create := &PatientCreate{}
				err = json.Unmarshal(body, create)
				assert.NoError(err)

				assert.Equal(testCase.create, create)

				b, err := json.Marshal(&Patient{})
				assert.NoError(err)

				w.Header().Set("Content-Type", "application/json")
				//nolint
				w.Write(b)
			}))
			defer srv.Close()

			client := NewHTTPClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
			svc := PatientService{client}

			created, res, err := svc.Create(context.Background(), testCase.create)
			assert.NotNil(created)
			assert.NotNil(res)
			assert.NoError(err)
		})
	}
}

func TestPatientService_Find(t *testing.T) {
	assert := assert.New(t)

	opts := &FindPatientsOptions{
		Pagination: &Pagination{
			Limit:  1,
			Offset: 2,
		},

		FirstName:        "first name",
		LastName:         "last name",
		DOB:              "dob",
		Sex:              "sex",
		InsuranceCompany: "insurance company",
		InsurancePlan:    "insurance plan",
		GroupID:          1,
		MemberID:         0,
		MasterPatient:    3,
		Practice:         4,
		LastModifiedGT:   time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
		LastModifiedGTE:  time.Date(2023, 5, 20, 0, 0, 0, 0, time.UTC),
		LastModifiedLT:   time.Date(2023, 5, 25, 0, 0, 0, 0, time.UTC),
		LastModifiedLTE:  time.Date(2023, 5, 30, 0, 0, 0, 0, time.UTC),
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/patients", r.URL.Path)

		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")
		dob := r.URL.Query().Get("dob")
		sex := r.URL.Query().Get("sex")
		insuranceCompany := r.URL.Query().Get("insurance_company")
		insurancePlan := r.URL.Query().Get("insurance_plan")
		groupID := r.URL.Query().Get("group_id")
		memberID := r.URL.Query().Get("member_id")
		masterPatient := r.URL.Query().Get("master_patient")
		practice := r.URL.Query().Get("practice")
		lastModifiedGT := r.URL.Query().Get("last_modified__gt")
		lastModifiedGTE := r.URL.Query().Get("last_modified__gte")
		lastModifiedLT := r.URL.Query().Get("last_modified__lt")
		lastModifiedLTE := r.URL.Query().Get("last_modified__lte")

		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")

		assert.Equal(opts.FirstName, firstName)
		assert.Equal(opts.LastName, lastName)
		assert.Equal(opts.DOB, dob)
		assert.Equal(opts.Sex, sex)
		assert.Equal(opts.InsuranceCompany, insuranceCompany)
		assert.Equal(opts.InsurancePlan, insurancePlan)
		assert.Equal(opts.GroupID, strToInt64(groupID))
		assert.Equal(opts.MemberID, strToInt64(memberID))
		assert.Equal(opts.MasterPatient, strToInt64(masterPatient))
		assert.Equal(opts.Practice, strToInt64(practice))
		assert.Equal(opts.LastModifiedGT.Format(time.RFC3339), lastModifiedGT)
		assert.Equal(opts.LastModifiedGTE.Format(time.RFC3339), lastModifiedGTE)
		assert.Equal(opts.LastModifiedLT.Format(time.RFC3339), lastModifiedLT)
		assert.Equal(opts.LastModifiedLTE.Format(time.RFC3339), lastModifiedLTE)

		assert.Equal(opts.Pagination.Limit, strToInt(limit))
		assert.Equal(opts.Pagination.Offset, strToInt(offset))

		b, err := json.Marshal(Response[[]*Patient]{
			Results: []*Patient{
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
	svc := PatientService{client}

	found, res, err := svc.Find(context.Background(), opts)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestPatientService_Get(t *testing.T) {
	assert := assert.New(t)

	var id int64 = 1

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/patients/"+strconv.FormatInt(id, 10), r.URL.Path)

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
	svc := PatientService{client}

	found, res, err := svc.Get(context.Background(), id)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestPatientService_Update(t *testing.T) {
	assert := assert.New(t)

	var id int64 = 1
	expected := &PatientUpdate{
		ActualName: Ptr("actual name"),
		Address: &PatientAddress{
			AddressLine1: "address line 1",
			AddressLine2: "address line 2",
			City:         "city",
			State:        "state",
			Zip:          "zip",
		},
		Consents: Ptr([]*PatientConsent{
			{
				ConsentType: "consent type",
				Expiration:  "expiration",
			},
		}),
		DOB: Ptr("dob"),
		Emails: Ptr([]*PatientEmail{
			{
				Email: "email",
			},
		}),
		Ethnicity:      Ptr("ethnicity"),
		FirstName:      Ptr("first name"),
		GenderIdentity: Ptr("gender identity name"),
		Insurances: Ptr([]*PatientInsuranceUpdate{
			{
				InsuranceCompany:       Ptr[int64](1),
				InsurancePlan:          Ptr[int64](2),
				Rank:                   "rank",
				Carrier:                Ptr("carrier"),
				MemberID:               Ptr("member ID"),
				GroupID:                Ptr("group ID"),
				Plan:                   Ptr("plan"),
				Phone:                  Ptr("phone"),
				Extension:              Ptr("extension"),
				Address:                Ptr("address"),
				Suite:                  Ptr("suite"),
				City:                   Ptr("city"),
				State:                  Ptr("state"),
				Zip:                    Ptr("zip"),
				Copay:                  Ptr("co-pay"),
				Deductible:             Ptr("deductable"),
				PaymentProgram:         Ptr("payment program"),
				InsuredPersonFirstName: Ptr("insured person first name"),
				InsuredPersonLastName:  Ptr("insured person last name"),
				InsuredPersonAddress:   Ptr("insured person address"),
				InsuredPersonCity:      Ptr("insured person city"),
				InsuredPersonState:     Ptr("insured person state"),
				InsuredPersonZip:       Ptr("insured person zip"),
				InsuredPersonID:        Ptr("insured person id"),
				InsuredPersonDOB:       Ptr("insured person DOB"),
				InsuredPersonGender:    Ptr("insured person gender"),
				InsuredPersonSSN:       Ptr("insured person SSN"),
				RelationshipToInsured:  Ptr("relationship to insured"),
				StartDate:              &civil.Date{Year: 2024, Month: 5, Day: 15},
				EndDate:                &civil.Date{Year: 3000, Month: 1, Day: 1},
			},
			{},
		}),
		LastName:          Ptr("last name"),
		LegalGenderMarker: Ptr("legal gender marker"),
		Metadata: &Metadata{
			Data: Ptr(map[string]string{
				"foo": "bar",
			}),
			ObjectID:      Ptr("object-id"),
			ObjectWebLink: Ptr("object-web-link"),
		},
		MiddleName: Ptr("middle name"),
		Notes:      Ptr("notes"),
		PatientStatus: &PatientStatusUpdate{
			InactiveReason: Ptr("other"),
			Status:         Ptr("inactive"),
		},
		Phones: Ptr([]*PatientPhone{
			{
				Phone:     "phone",
				PhoneType: "phone type",
			},
		}),
		PreferredLanguage:      Ptr("preferred language"),
		PrimaryCareProviderNPI: Ptr("primary care provider NPI"),
		PrimaryPhysician:       Ptr[int64](1),
		Pronouns:               Ptr("pronouns"),
		Race:                   Ptr("race"),
		Sex:                    Ptr("sex"),
		SexualOrientation:      Ptr("sexual orientation"),
		SSN:                    Ptr("ssn"),
		Tags:                   []string{"PAID", "Test Patient"},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodPatch, r.Method)
		assert.Equal("/patients/"+strconv.FormatInt(id, 10), r.URL.Path)

		body, err := io.ReadAll(r.Body)
		assert.NoError(err)

		actual := &PatientUpdate{}
		err = json.Unmarshal(body, actual)
		assert.NoError(err)

		assert.Equal(expected, actual)

		b, err := json.Marshal(&Patient{})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := PatientService{client}

	updated, res, err := svc.Update(context.Background(), id, expected)
	assert.NotNil(updated)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestPatientService_Update_empty_arrays(t *testing.T) {
	testCases := map[string]struct {
		consents   *[]*PatientConsent
		emails     *[]*PatientEmail
		insurances *[]*PatientInsuranceUpdate
		phones     *[]*PatientPhone
	}{
		"pointers to nil slices": {
			consents:   Ptr([]*PatientConsent(nil)),
			emails:     Ptr([]*PatientEmail(nil)),
			insurances: Ptr([]*PatientInsuranceUpdate(nil)),
			phones:     Ptr([]*PatientPhone(nil)),
		},
		"pointers to non-nil, empty slices": {
			consents:   &[]*PatientConsent{},
			emails:     &[]*PatientEmail{},
			insurances: &[]*PatientInsuranceUpdate{},
			phones:     &[]*PatientPhone{},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			var id int64 = 1
			expected := &PatientUpdate{
				ActualName: Ptr("actual name"),
				Address: &PatientAddress{
					AddressLine1: "address line 1",
					AddressLine2: "address line 2",
					City:         "city",
					State:        "state",
					Zip:          "zip",
				},
				Consents:          testCase.consents,
				DOB:               Ptr("dob"),
				Emails:            testCase.emails,
				Ethnicity:         Ptr("ethnicity"),
				FirstName:         Ptr("first name"),
				GenderIdentity:    Ptr("gender identity name"),
				Insurances:        testCase.insurances,
				LastName:          Ptr("last name"),
				LegalGenderMarker: Ptr("legal gender marker"),
				MiddleName:        Ptr("middle name"),
				Notes:             Ptr("notes"),
				PatientStatus: &PatientStatusUpdate{
					InactiveReason: Ptr("other"),
					Status:         Ptr("inactive"),
				},
				Phones:                 testCase.phones,
				PreferredLanguage:      Ptr("preferred language"),
				PrimaryCareProviderNPI: Ptr("primary care provider NPI"),
				PrimaryPhysician:       Ptr[int64](1),
				Pronouns:               Ptr("pronouns"),
				Race:                   Ptr("race"),
				Sex:                    Ptr("sex"),
				SexualOrientation:      Ptr("sexual orientation"),
				SSN:                    Ptr("ssn"),
			}

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tokenRequest(w, r) {
					return
				}

				assert.Equal(http.MethodPatch, r.Method)
				assert.Equal("/patients/"+strconv.FormatInt(id, 10), r.URL.Path)

				body, err := io.ReadAll(r.Body)
				assert.NoError(err)

				actual := &PatientUpdate{}
				err = json.Unmarshal(body, actual)
				assert.NoError(err)

				assert.Empty(cmp.Diff(
					expected,
					actual,
					// NonNullJSONArray always marshals empty JSON arrays to nil slices, so we want to equate nil and
					// empty slices for the "pointers to non-nil, empty slices" test case.
					cmpopts.EquateEmpty(),
				))

				assert.Contains(string(body), `"consents":[]`)
				assert.Contains(string(body), `"emails":[]`)
				assert.Contains(string(body), `"insurances":[]`)
				assert.Contains(string(body), `"phones":[]`)

				b, err := json.Marshal(&Patient{})
				assert.NoError(err)

				w.Header().Set("Content-Type", "application/json")
				//nolint
				w.Write(b)
			}))
			defer srv.Close()

			client := NewHTTPClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
			svc := PatientService{client}

			updated, res, err := svc.Update(context.Background(), id, expected)
			assert.NotNil(updated)
			assert.NotNil(res)
			assert.NoError(err)
		})
	}
}

func TestPatientService_Update_omitted_array_keys(t *testing.T) {
	assert := assert.New(t)

	var id int64 = 1
	expected := &PatientUpdate{
		ActualName: Ptr("actual name"),
		Address: &PatientAddress{
			AddressLine1: "address line 1",
			AddressLine2: "address line 2",
			City:         "city",
			State:        "state",
			Zip:          "zip",
		},
		DOB:               Ptr("dob"),
		Ethnicity:         Ptr("ethnicity"),
		FirstName:         Ptr("first name"),
		GenderIdentity:    Ptr("gender identity name"),
		LastName:          Ptr("last name"),
		LegalGenderMarker: Ptr("legal gender marker"),
		MiddleName:        Ptr("middle name"),
		Notes:             Ptr("notes"),
		PatientStatus: &PatientStatusUpdate{
			InactiveReason: Ptr("other"),
			Status:         Ptr("inactive"),
		},
		PreferredLanguage:      Ptr("preferred language"),
		PrimaryCareProviderNPI: Ptr("primary care provider NPI"),
		PrimaryPhysician:       Ptr[int64](1),
		Pronouns:               Ptr("pronouns"),
		Race:                   Ptr("race"),
		Sex:                    Ptr("sex"),
		SexualOrientation:      Ptr("sexual orientation"),
		SSN:                    Ptr("ssn"),
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodPatch, r.Method)
		assert.Equal("/patients/"+strconv.FormatInt(id, 10), r.URL.Path)

		body, err := io.ReadAll(r.Body)
		assert.NoError(err)

		actual := &PatientUpdate{}
		err = json.Unmarshal(body, actual)
		assert.NoError(err)

		assert.Equal(expected, actual)

		assert.NotContains(string(body), `"consents"`)
		assert.NotContains(string(body), `"emails"`)
		assert.NotContains(string(body), `"insurances"`)
		assert.NotContains(string(body), `"phones"`)

		b, err := json.Marshal(&Patient{})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := PatientService{client}

	updated, res, err := svc.Update(context.Background(), id, expected)
	assert.NotNil(updated)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestPatientService_Delete(t *testing.T) {
	assert := assert.New(t)

	var id int64 = 1

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodDelete, r.Method)
		assert.Equal("/patients/"+strconv.FormatInt(id, 10), r.URL.Path)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := PatientService{client}

	res, err := svc.Delete(context.Background(), id)
	assert.NotNil(res)
	assert.NoError(err)
}

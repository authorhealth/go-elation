package elation

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type PatientServicer interface {
	Create(ctx context.Context, create *PatientCreate) (*Patient, *http.Response, error)
	Find(ctx context.Context, opts *FindPatientsOptions) (*Response[[]*Patient], *http.Response, error)
	Get(ctx context.Context, id int64) (*Patient, *http.Response, error)
	Update(ctx context.Context, id int64, update *PatientUpdate) (*Patient, *http.Response, error)
	Delete(ctx context.Context, id int64) (*http.Response, error)
}

var _ PatientServicer = (*PatientService)(nil)

type PatientService struct {
	client *Client
}

type PatientCreate struct {
	LastName          string          `json:"last_name"`
	FirstName         string          `json:"first_name"`
	Sex               string          `json:"sex"`
	DOB               string          `json:"dob"`
	PrimaryPhysician  int64           `json:"primary_physician"`
	CaregiverPractice int64           `json:"caregiver_practice"`
	Address           *PatientAddress `json:"address,omitempty"`
	Phones            []*PatientPhone `json:"phones,omitempty"`
	Emails            []*PatientEmail `json:"emails,omitempty"`
}

func (s *PatientService) Create(ctx context.Context, create *PatientCreate) (*Patient, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "create patient", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &Patient{}

	res, err := s.client.request(ctx, http.MethodPost, "/patients", nil, create, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

type Patient struct {
	ID                     int64               `json:"id"`
	FirstName              string              `json:"first_name"`
	MiddleName             string              `json:"middle_name"`
	LastName               string              `json:"last_name"`
	ActualName             string              `json:"actual_name"`
	GenderIdentity         string              `json:"gender_identity"`
	LegalGenderMarker      string              `json:"legal_gender_marker"`
	Pronouns               string              `json:"pronouns"`
	Sex                    string              `json:"sex"`
	SexualOrientation      string              `json:"sexual_orientation"`
	PrimaryPhysician       int64               `json:"primary_physician"`
	CaregiverPractice      int64               `json:"caregiver_practice"`
	DOB                    string              `json:"dob"`
	SSN                    string              `json:"ssn"`
	Race                   string              `json:"race"`
	PreferredLanguage      string              `json:"preferred_language"`
	Ethnicity              string              `json:"ethnicity"`
	Notes                  string              `json:"notes"`
	VIP                    bool                `json:"vip"`
	Address                *PatientAddress     `json:"address"`
	Phones                 []*PatientPhone     `json:"phones"`
	Emails                 []*PatientEmail     `json:"emails"`
	Guarantor              *PatientGuarantor   `json:"guarantor"`
	Insurances             []*PatientInsurance `json:"insurances"`
	DeletedInsurances      []*PatientInsurance `json:"deleted_insurances"`
	Tags                   []string            `json:"tags"`
	PatientStatus          *PatientStatus      `json:"patient_status"`
	Preference             *PatientPreference  `json:"preference"`
	EmergencyContact       *PatientContact     `json:"emergency_contact"`
	PrimaryCareProvider    int64               `json:"primary_care_provider"`
	PrimaryCareProviderNPI string              `json:"primary_care_provider_npi"`
	PreviousFirstName      string              `json:"previous_first_name"`
	PreviousLastName       string              `json:"previous_last_name"`
	MasterPatient          *int64              `json:"master_patient"`
	Employer               *PatientEmployer    `json:"employer"`
	Consents               []*PatientConsent   `json:"consents"`
	Metadata               any                 `json:"metadata"`
	CreatedDate            time.Time           `json:"created_date"`
	DeletedDate            *time.Time          `json:"deleted_date"`
	MergedIntoChart        int64               `json:"merged_into_chart"`
}

type PatientAddress struct {
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city"`
	State        string `json:"state"`
	Zip          string `json:"zip"`
}

type PatientPhone struct {
	Phone       string     `json:"phone"`
	PhoneType   string     `json:"phone_type"`
	CreatedDate time.Time  `json:"created_date"`
	DeletedDate *time.Time `json:"deleted_date"`
}

type PatientEmail struct {
	Email       string     `json:"email"`
	CreatedDate time.Time  `json:"created_date"`
	DeletedDate *time.Time `json:"deleted_date"`
}

type PatientGuarantor struct {
	ID           int64  `json:"id"`
	Address      string `json:"address"`
	City         string `json:"city"`
	State        string `json:"state"`
	Zip          string `json:"zip"`
	Phone        string `json:"phone"`
	Relationship string `json:"relationship"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	MiddleName   string `json:"middle_name"`
}

type PatientInsurance struct {
	ID                     int64      `json:"id"`
	InsuranceCompany       int64      `json:"insurance_company,omitempty"`
	InsurancePlan          int64      `json:"insurance_plan,omitempty"`
	Rank                   string     `json:"rank"`
	Carrier                string     `json:"carrier"`
	MemberID               string     `json:"member_id"`
	GroupID                string     `json:"group_id"`
	Plan                   string     `json:"plan"`
	Phone                  string     `json:"phone"`
	Extension              string     `json:"extension"`
	Address                string     `json:"address"`
	Suite                  string     `json:"suite"`
	City                   string     `json:"city"`
	State                  string     `json:"state"`
	Zip                    string     `json:"zip"`
	Copay                  any        `json:"copay"`
	Deductible             any        `json:"deductible"`
	PaymentProgram         string     `json:"payment_program"`
	InsuredPersonFirstName string     `json:"insured_person_first_name"`
	InsuredPersonLastName  string     `json:"insured_person_last_name"`
	InsuredPersonAddress   string     `json:"insured_person_address"`
	InsuredPersonCity      string     `json:"insured_person_city"`
	InsuredPersonState     string     `json:"insured_person_state"`
	InsuredPersonZip       string     `json:"insured_person_zip"`
	InsuredPersonID        string     `json:"insured_person_id"`
	InsuredPersonDOB       string     `json:"insured_person_dob"`
	InsuredPersonGender    string     `json:"insured_person_gender"`
	InsuredPersonSSN       string     `json:"insured_person_ssn"`
	RelationshipToInsured  string     `json:"relationship_to_insured"`
	CreatedDate            time.Time  `json:"created_date"`
	DeletedDate            *time.Time `json:"deleted_date"`
	StartDate              string     `json:"start_date,omitempty"`
	EndDate                string     `json:"end_date,omitempty"`
}

type PatientStatus struct {
	DeceasedDate     string    `json:"deceased_date"`
	InactiveReason   string    `json:"inactive_reason"`
	LastStatusChange time.Time `json:"last_status_change"`
	Notes            string    `json:"notes"`
	Status           string    `json:"status"`
}

type PatientPreference struct {
	PreferredPharmacy1 any `json:"preferred_pharmacy_1"`
	PreferredPharmacy2 any `json:"preferred_pharmacy_2"`
}

type PatientContact struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Relationship string `json:"relationship"`
	Phone        string `json:"phone"`
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city"`
	State        string `json:"state"`
	Zip          string `json:"zip"`
}

type PatientEmployer struct {
	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	Zip        string `json:"zip"`
	EmployerID string `json:"employer_id"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
}

type PatientConsent struct {
	ID          int64  `json:"id"`
	ConsentType string `json:"consent_type"`
	Expiration  string `json:"expiration"`
}

type FindPatientsOptions struct {
	*Pagination

	FirstName        string    `url:"first_name,omitempty"`
	LastName         string    `url:"last_name,omitempty"`
	DOB              string    `url:"dob,omitempty"`
	Sex              string    `url:"sex,omitempty"`
	InsuranceCompany string    `url:"insurance_company,omitempty"`
	InsurancePlan    string    `url:"insurance_plan,omitempty"`
	GroupID          int64     `url:"group_id,omitempty"`
	MemberID         int64     `url:"member_id,omitempty"`
	MasterPatient    int64     `url:"master_patient,omitempty"`
	Practice         int64     `url:"practice,omitempty"`
	LastModifiedGT   time.Time `url:"last_modified_gt,omitempty"`
	LastModifiedGTE  time.Time `url:"last_modified_gte,omitempty"`
	LastModifiedLT   time.Time `url:"last_modified_lt,omitempty"`
	LastModifiedLTE  time.Time `url:"last_modified_lte,omitempty"`
}

func (s *PatientService) Find(ctx context.Context, opts *FindPatientsOptions) (*Response[[]*Patient], *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "find patients", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &Response[[]*Patient]{}

	res, err := s.client.request(ctx, http.MethodGet, "/patients", opts, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *PatientService) Get(ctx context.Context, id int64) (*Patient, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "get patient", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.patient_id", id)))
	defer span.End()

	out := &Patient{}

	res, err := s.client.request(ctx, http.MethodGet, "/patients/"+strconv.FormatInt(id, 10), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

type PatientUpdate struct {
	ActualName             *string              `json:"actual_name,omitempty"`
	Address                *PatientAddress      `json:"address,omitempty"`
	Consents               []*PatientConsent    `json:"consents,omitempty"`
	DOB                    *string              `json:"dob,omitempty"`
	Emails                 []*PatientEmail      `json:"emails,omitempty"`
	Ethnicity              *string              `json:"ethnicity,omitempty"`
	FirstName              *string              `json:"first_name,omitempty"`
	GenderIdentity         *string              `json:"gender_identity,omitempty"`
	Insurances             []*PatientInsurance  `json:"insurances,omitempty"`
	LastName               *string              `json:"last_name,omitempty"`
	LegalGenderMarker      *string              `json:"legal_gender_marker,omitempty"`
	MiddleName             *string              `json:"middle_name,omitempty"`
	Notes                  *string              `json:"notes,omitempty"`
	PatientStatus          *PatientStatusUpdate `json:"patient_status,omitempty"`
	Phones                 []*PatientPhone      `json:"phones,omitempty"`
	PreferredLanguage      *string              `json:"preferred_language,omitempty"`
	PrimaryCareProviderNPI *string              `json:"primary_care_provider_npi,omitempty"`
	PrimaryPhysician       *int64               `json:"primary_physician,omitempty"`
	Pronouns               *string              `json:"pronouns,omitempty"`
	Race                   *string              `json:"race,omitempty"`
	Sex                    *string              `json:"sex,omitempty"`
	SexualOrientation      *string              `json:"sexual_orientation,omitempty"`
	SSN                    *string              `json:"ssn,omitempty"`
}

type PatientStatusUpdate struct {
	InactiveReason *string `json:"inactive_reason,omitempty"`
	Status         *string `json:"status,omitempty"`
}

func (s *PatientService) Update(ctx context.Context, id int64, update *PatientUpdate) (*Patient, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "update patient", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.patient_id", id)))
	defer span.End()

	out := &Patient{}

	res, err := s.client.request(ctx, http.MethodPatch, "/patients/"+strconv.FormatInt(id, 10), nil, update, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *PatientService) Delete(ctx context.Context, id int64) (*http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "delete patient", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.patient_id", id)))
	defer span.End()

	res, err := s.client.request(ctx, http.MethodDelete, "/patients/"+strconv.FormatInt(id, 10), nil, nil, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return res, fmt.Errorf("making request: %w", err)
	}

	return res, nil
}

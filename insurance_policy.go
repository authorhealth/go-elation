package elation

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type InsurancePolicyServicer interface {
	Create(ctx context.Context, patientID int64, create *InsurancePolicyCreate) (*InsurancePolicy, *http.Response, error)
	Find(ctx context.Context, patientID int64, opts *FindInsurancePoliciesOptions) (*InsurancePolicyResponse, *http.Response, error)
	Get(ctx context.Context, patientID int64, id int64) (*InsurancePolicy, *http.Response, error)
	Update(ctx context.Context, patientID int64, id int64, update *InsurancePolicy) (*InsurancePolicy, *http.Response, error)
	Delete(ctx context.Context, patientID int64, id int64) (*http.Response, error)
}

var _ InsurancePolicyServicer = (*InsurancePolicyService)(nil)

type InsurancePolicyService struct {
	client *HTTPClient
}

type InsurancePolicyResponse struct {
	// TODO: Determine how the pagination fields here work / if they work yet

	Results []*InsurancePolicy `json:"results"`
}

type InsurancePolicy struct {
	ID int64 `json:"id"`

	PracticeID       int64  `json:"practice_id"`
	PatientID        int64  `json:"patient_id"`
	PatientFirstName string `json:"patient_first_name"`
	PatientLastName  string `json:"patient_last_name"`
	PatientDOB       string `json:"patient_dob"` // Format: YYYY-MM-DD

	Status      string  `json:"status"` // enum: "active", "inactive"
	Rank        *int64  `json:"rank"`   // enum: 1=primary, 2=secondary, 3=tertiary
	CarrierID   *int64  `json:"carrier_id"`
	CarrierName *string `json:"carrier_name"`
	PlanID      *int64  `json:"plan_id"`
	PlanName    *string `json:"plan_name"`
	GroupID     *string `json:"group_id"`
	MemberID    *string `json:"member_id"`
	Copay       *string `json:"copay"`
	Deductible  *string `json:"deductible"`
	StartDate   *string `json:"start_date"` // Format: YYYY-MM-DD
	EndDate     *string `json:"end_date"`   // Format: YYYY-MM-DD

	Phone      *string `json:"phone"`       // Carrier phone
	Extension  *string `json:"extension"`   // Carrier phone extension
	Address    *string `json:"address"`     // Carrier address
	Suite      *string `json:"suite"`       // Carrier suite
	City       *string `json:"city"`        // Carrier city
	State      *string `json:"state"`       // Carrier state
	Zip        *string `json:"zip"`         // Carrier zip code
	CountyCode *string `json:"county_code"` // Carrier phone country code

	InsuredPersonFirstName  *string `json:"insured_person_first_name"`
	InsuredPersonLastName   *string `json:"insured_person_last_name"`
	InsuredPersonAddress    *string `json:"insured_person_address"`
	InsuredPersonCity       *string `json:"insured_person_city"`
	InsuredPersonState      *string `json:"insured_person_state"`
	InsuredPersonZip        *string `json:"insured_person_zip"`
	InsuredPersonID         *string `json:"insured_person_id"`
	InsuredPersonDOB        *string `json:"insured_person_dob"` // date format: YYYY-MM-DD
	InsuredPersonSexAtBirth *string `json:"insured_person_sex_at_birth"`
	InsuredPersonSSN        *string `json:"insured_person_ssn"`
	RelationshipToInsured   *string `json:"relationship_to_insured"` // enum: "self", "child", "spouse", "other", default: "self"

	PaymentProgram        *string `json:"payment_program"` // enum: "medicare_part_b", "medicare_advantage", "medicaid", "commercial_hmsa", "commercial_sfhp", "commercial_other", "workers_compensation"
	MedicareSecondaryCode *string `json:"medicare_secondary_code"`
}

type InsurancePolicyCreate struct {
	PracticeID       int64  `json:"practice_id"`
	PatientID        int64  `json:"patient_id"`
	PatientFirstName string `json:"patient_first_name"`
	PatientLastName  string `json:"patient_last_name"`
	PatientDOB       string `json:"patient_dob"` // Format: YYYY-MM-DD
	Status           string `json:"status"`      // enum: "active", "inactive"

	Rank                    *int64  `json:"rank,omitempty"` // enum: 1=primary, 2=secondary, 3=tertiary
	CarrierID               *int64  `json:"carrier_id,omitempty"`
	CarrierName             *string `json:"carrier_name,omitempty"`
	PlanID                  *int64  `json:"plan_id,omitempty"`
	PlanName                *string `json:"plan_name,omitempty"`
	GroupID                 *string `json:"group_id,omitempty"`
	MemberID                *string `json:"member_id,omitempty"`
	Copay                   *string `json:"copay,omitempty"`
	Deductible              *string `json:"deductible,omitempty"`
	StartDate               *string `json:"start_date,omitempty"`  // Format: YYYY-MM-DD
	EndDate                 *string `json:"end_date,omitempty"`    // Format: YYYY-MM-DD
	Phone                   *string `json:"phone,omitempty"`       // Carrier phone
	Extension               *string `json:"extension,omitempty"`   // Carrier phone extension
	Address                 *string `json:"address,omitempty"`     // Carrier address
	Suite                   *string `json:"suite,omitempty"`       // Carrier suite
	City                    *string `json:"city,omitempty"`        // Carrier city
	State                   *string `json:"state,omitempty"`       // Carrier state
	Zip                     *string `json:"zip,omitempty"`         // Carrier zip code
	CountyCode              *string `json:"county_code,omitempty"` // Carrier phone country code
	InsuredPersonFirstName  *string `json:"insured_person_first_name,omitempty"`
	InsuredPersonLastName   *string `json:"insured_person_last_name,omitempty"`
	InsuredPersonAddress    *string `json:"insured_person_address,omitempty"`
	InsuredPersonCity       *string `json:"insured_person_city,omitempty"`
	InsuredPersonState      *string `json:"insured_person_state,omitempty"`
	InsuredPersonZip        *string `json:"insured_person_zip,omitempty"`
	InsuredPersonID         *string `json:"insured_person_id,omitempty"`
	InsuredPersonDOB        *string `json:"insured_person_dob,omitempty"`
	InsuredPersonSexAtBirth *string `json:"insured_person_sex_at_birth,omitempty"`
	InsuredPersonSSN        *string `json:"insured_person_ssn,omitempty"`
	RelationshipToInsured   *string `json:"relationship_to_insured,omitempty"`
	PaymentProgram          *string `json:"payment_program,omitempty"`
	MedicareSecondaryCode   *string `json:"medicare_secondary_code,omitempty"`
}

type FindInsurancePoliciesOptions struct {
	// TODO: Understand how pagination functions for this endpoint

	ActiveOnly *bool `url:"active_only,omitempty"`
}

func (s *InsurancePolicyService) Create(ctx context.Context, patientID int64, create *InsurancePolicyCreate) (*InsurancePolicy, *http.Response, error) {
	ctx, span := s.client.tracer.Start(
		ctx,
		"create insurance policy",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attribute.Int64("elation.patient_id", patientID)),
	)
	defer span.End()

	out := &InsurancePolicy{}

	res, err := s.client.request(ctx, http.MethodPost, "/patients/"+strconv.FormatInt(patientID, 10)+"/policies", nil, create, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *InsurancePolicyService) Find(ctx context.Context, patientID int64, opts *FindInsurancePoliciesOptions) (*InsurancePolicyResponse, *http.Response, error) {
	ctx, span := s.client.tracer.Start(
		ctx,
		"find insurance policies",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attribute.Int64("elation.patient_id", patientID)),
	)
	defer span.End()

	out := &InsurancePolicyResponse{}

	res, err := s.client.request(ctx, http.MethodGet, "/patients/"+strconv.FormatInt(patientID, 10)+"/policies", opts, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *InsurancePolicyService) Get(ctx context.Context, patientID int64, id int64) (*InsurancePolicy, *http.Response, error) {
	ctx, span := s.client.tracer.Start(
		ctx,
		"get insurance policy",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attribute.Int64("elation.patient_id", patientID)),
		trace.WithAttributes(attribute.Int64("elation.insurance_policy_id", id)),
	)
	defer span.End()

	out := &InsurancePolicy{}

	res, err := s.client.request(ctx, http.MethodGet, "/patients/"+strconv.FormatInt(patientID, 10)+"/policies/"+strconv.FormatInt(id, 10), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *InsurancePolicyService) Update(ctx context.Context, patientID int64, id int64, update *InsurancePolicy) (*InsurancePolicy, *http.Response, error) {
	ctx, span := s.client.tracer.Start(
		ctx,
		"update insurance policy",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attribute.Int64("elation.patient_id", patientID)),
		trace.WithAttributes(attribute.Int64("elation.insurance_policy_id", id)),
	)
	defer span.End()

	out := &InsurancePolicy{}

	res, err := s.client.request(ctx, http.MethodPut, "/patients/"+strconv.FormatInt(patientID, 10)+"/policies/"+strconv.FormatInt(id, 10), nil, update, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *InsurancePolicyService) Delete(ctx context.Context, patientID int64, id int64) (*http.Response, error) {
	ctx, span := s.client.tracer.Start(
		ctx,
		"delete insurance policy",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attribute.Int64("elation.patient_id", patientID)),
		trace.WithAttributes(attribute.Int64("elation.insurance_policy_id", id)),
	)
	defer span.End()

	res, err := s.client.request(ctx, http.MethodDelete, "/patients/"+strconv.FormatInt(patientID, 10)+"/policies/"+strconv.FormatInt(id, 10), nil, nil, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return res, fmt.Errorf("making request: %w", err)
	}

	return res, nil
}

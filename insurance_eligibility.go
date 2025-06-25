package elation

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type InsuranceEligibilityServicer interface {
	Create(ctx context.Context, patientInsuranceID int64, create *InsuranceEligibilityCreate) (*InsuranceEligibility, *http.Response, error)
	Get(ctx context.Context, patientInsuranceID int64) (*InsuranceEligibility, *http.Response, error)
	GetFullReport(ctx context.Context, patientInsuranceID int64) (*InsuranceEligibilityFullReport, *http.Response, error)
}

var _ InsuranceEligibilityServicer = (*InsuranceEligibilityService)(nil)

type InsuranceEligibilityService struct {
	client *HTTPClient
}

type InsuranceEligibilityCreate struct {
	GroupID   int64  `json:"group_id,omitempty"`
	GroupName string `json:"group_name,omitempty"`
	GroupNPI  string `json:"group_npi,omitempty"`
}

func (s *InsuranceEligibilityService) Create(ctx context.Context, patientInsuranceID int64, create *InsuranceEligibilityCreate) (*InsuranceEligibility, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "create insurance eligibility", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.patient_insurance_id", patientInsuranceID)))
	defer span.End()

	out := &InsuranceEligibility{}

	res, err := s.client.request(ctx, http.MethodPost, fmt.Sprintf("/patient_insurances/%d/eligibility/", patientInsuranceID), nil, create, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

type InsuranceEligibility struct {
	EligibilityDetails        *InsuranceEligibilityDetails `json:"eligibility_details"`
	EligibilityCheckTimestamp timeWithOptionalZone         `json:"eligibility_check_timestamp"`
	EligibilityStatus         string                       `json:"eligibility_status"`
	PatientID                 int64                        `json:"patient_id"`
	PatientInsuranceID        int64                        `json:"patient_insurance_id"`
	PracticeID                int64                        `json:"practice_id"`
}

type InsuranceEligibilityDetails struct {
	Copay                        any    `json:"copay"`
	CopayAmbiguous               bool   `json:"copay_ambiguous"`
	Coinsurance                  string `json:"coinsurance"`
	CoinsuranceAmbiguous         bool   `json:"coinsurance_ambiguous"`
	Deductible                   string `json:"deductible"`
	DeductibleAmbiguous          bool   `json:"deductible_ambiguous"`
	DeductibleRemaining          string `json:"deductible_remaining"`
	DedictibleRemainingAmbiguous bool   `json:"deductible_remaining_ambiguous"`
	Errors                       []any  `json:"errors"`
}

func (s *InsuranceEligibilityService) Get(ctx context.Context, patientInsuranceID int64) (*InsuranceEligibility, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "get insurance eligibility", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.patient_insurance_id", patientInsuranceID)))
	defer span.End()

	out := &InsuranceEligibility{}

	res, err := s.client.request(ctx, http.MethodGet, fmt.Sprintf("/patient_insurances/%d/eligibility/", patientInsuranceID), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

type InsuranceEligibilityFullReport struct {
	EligibilityCheckTimestamp timeWithOptionalZone    `json:"eligibility_check_timestamp"`
	PatientID                 int64                   `json:"patient_id"`
	PatientInsuranceID        int64                   `json:"patient_insurance_id"`
	PracticeID                int64                   `json:"practice_id"`
	ServiceTypeCode           string                  `json:"service_type_code"`
	Subscriber                IEFRSubscriber          `json:"subscriber"`
	EligibilityProvider       IEFREligibilityProvider `json:"eligibility_provider"`
	PlanDetails               IEFRPlanDetails         `json:"plan_details"`
	Benefits                  IEFRBenefits            `json:"benefits"`
}

type IEFRSubscriber struct {
	Address    IEFRSubscriberAddress `json:"address"`
	FirstName  string                `json:"first_name"`
	LastName   string                `json:"last_name"`
	MemberID   string                `json:"member_id"`
	MiddleName string                `json:"middle_name"`
	Name       string                `json:"name"`
	SexAtBirth string                `json:"sex_at_birth"`
}

type IEFRSubscriberAddress struct {
	Address1   string `json:"address_1"`
	Address2   string `json:"address_2"`
	City       string `json:"city"`
	PostalCode string `json:"postal_code"`
	State      string `json:"state"`
}

type IEFREligibilityProvider struct {
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	NPI              string `json:"npi"`
	OrganizationName string `json:"organization_name"`
}

type IEFRPlanDetails struct {
	Carrier       string                  `json:"carrier"`
	EndDate       string                  `json:"end_date"`
	GroupName     string                  `json:"group_name"`
	GroupNumber   string                  `json:"group_number"`
	InsuranceType string                  `json:"insurance_type"`
	MemberID      string                  `json:"member_id"`
	PlanName      string                  `json:"plan_name"`
	PlanNumber    string                  `json:"plan_number"`
	PolicyName    string                  `json:"policy_name"`
	PolicyNumber  string                  `json:"policy_number"`
	Provider      IEFRPlanDetailsProvider `json:"provider"`
	StartDate     string                  `json:"start_date"`
}

type IEFRPlanDetailsProvider struct {
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	NPI              string `json:"npi"`
	OrganizationName string `json:"organization_name"`
}

type IEFRBenefits struct {
	Copay       []map[string]any `json:"copay"`
	Coinsurance []any            `json:"coinsurance"`
	Deductible  []any            `json:"deductible"`
	OutOfPocket []any            `json:"out_of_pocket"`
	Limitations []any            `json:"limitations"`
}

func (s *InsuranceEligibilityService) GetFullReport(ctx context.Context, patientInsuranceID int64) (*InsuranceEligibilityFullReport, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "get insurance eligibility full report", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.patient_insurance_id", patientInsuranceID)))
	defer span.End()

	out := &InsuranceEligibilityFullReport{}

	res, err := s.client.request(ctx, http.MethodGet, fmt.Sprintf("/patient_insurances/%d/eligibility_full_report/", patientInsuranceID), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

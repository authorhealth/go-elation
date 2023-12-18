package elation

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type InsuranceEligibilityServicer interface {
	Create(ctx context.Context, id int64, create *InsuranceEligibilityCreate) (*InsuranceEligibility, *http.Response, error)
	Get(ctx context.Context, id int64) (*InsuranceEligibility, *http.Response, error)
	GetFullReport(ctx context.Context, id int64) (*InsuranceEligibilityFullReport, *http.Response, error)
}

var _ InsuranceEligibilityServicer = (*InsuranceEligibilityService)(nil)

type InsuranceEligibilityService struct {
	client *Client
}

type InsuranceEligibilityCreate struct {
	GroupID   int64  `json:"group_id,omitempty"`
	GroupName string  `json:"group_name,omitempty"`
	GroupNPI  string `json:"group_npi,omitempty"`
}

func (s *InsuranceEligibilityService) Create(ctx context.Context, id int64, create *InsuranceEligibilityCreate) (*InsuranceEligibility, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "create insurance eligibility", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.patient_insurance_id", id)))
	defer span.End()

	out := &InsuranceEligibility{}

	res, err := s.client.request(ctx, http.MethodPost, fmt.Sprintf("/patient_insurances/%d/eligibility/", id), nil, create, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

type InsuranceEligibility struct {
	EligibilityDetails        *InsuranceEligibilityDetails `json:"eligibility_details"`
	EligibilityCheckTimestamp time.Time                    `json:"eligibility_check_timestamp"`
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

func (s *InsuranceEligibilityService) Get(ctx context.Context, id int64) (*InsuranceEligibility, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "get insurance eligibility", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.patient_insurance_id", id)))
	defer span.End()

	out := &InsuranceEligibility{}

	res, err := s.client.request(ctx, http.MethodGet, fmt.Sprintf("/patient_insurances/%d/eligibility/", id), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

type InsuranceEligibilityFullReport struct {
	EligibilityCheckTimestamp time.Time               `json:"eligibility_check_timestamp"` //: "2022-04-04T01:52:38.662938+00:00",
	PatientID                 int64                   `json:"patient_id"`                  //: 123,
	PatientInsuranceID        int64                   `json:"patient_insurance_id"`        //: 456,
	PracticeID                int64                   `json:"practice_id"`                 //: 789,
	ServiceTypeCode           string                  `json:"service_type_code"`           //: "30",
	Subscriber                IEFRSubscriber          `json:"subscriber"`
	EligibilityProvider       IEFREligibilityProvider `json:"eligibility_provider"`
	PlanDetails               IEFRPlanDetails         `json:"plan_details"`
	Benefits                  IEFRBenefits            `json:"benefits"`
}

type IEFRSubscriber struct {
	Address    IEFRSubscriberAddress `json:"address"`
	FirstName  string                `json:"first_name"`   //: "Paula",
	LastName   string                `json:"last_name"`    //: "Patient",
	MemberID   string                `json:"member_id"`    //: "test_member_id",
	MiddleName string                `json:"middle_name"`  //: "P",
	Name       string                `json:"name"`         //: "Paula Patient",
	SexAtBirth string                `json:"sex_at_birth"` //: "F"
}

type IEFRSubscriberAddress struct {
	Address1   string `json:"address_1"`   //: "550 15th Street",
	Address2   string `json:"address_2"`   //: "21",
	City       string `json:"city"`        //: "San Francisco",
	PostalCode string `json:"postal_code"` //: "94103",
	State      string `json:"state"`       //: "CA"
}

type IEFREligibilityProvider struct {
	FirstName        string `json:"first_name"`        //: "Gary",
	LastName         string `json:"last_name"`         //: "Leung",
	NPI              string `json:"npi"`               //: "1234567890",
	OrganizationName string `json:"organization_name"` //: "Elation North"
}

type IEFRPlanDetails struct {
	Carrier       string                  `json:"carrier"`        //: "Blue Cross",
	EndDate       string                  `json:"end_date"`       //: null,
	GroupName     string                  `json:"group_name"`     //: "Elation Health, Inc",
	GroupNumber   string                  `json:"group_number"`   //: "test_group_number",
	InsuranceType string                  `json:"insurance_type"` //: null,
	MemberID      string                  `json:"member_id"`      //: "test_member_id",
	PlanName      string                  `json:"plan_name"`      //: "Blue Cross Blue Shield",
	PlanNumber    string                  `json:"plan_number"`    //: "test_plan_number",
	PolicyName    string                  `json:"policy_name"`    //: null,
	PolicyNumber  string                  `json:"policy_number"`  //: null,
	Provider      IEFRPlanDetailsProvider `json:"provider"`
	StartDate     string                  `json:"start_date"` //: "2022-04-04"
}

type IEFRPlanDetailsProvider struct {
	FirstName        string `json:"first_name"`        //: "Gary",
	LastName         string `json:"last_name"`         //: "Leung",
	NPI              string `json:"npi"`               //: "1234567890",
	OrganizationName string `json:"organization_name"` //: "Elation North"
}

type IEFRBenefits struct {
	Copay       []any `json:"copay"`         //: [],
	Coinsurance []any `json:"coinsurance"`   //: [],
	Deductible  []any `json:"deductible"`    //: [],
	OutOfPocket []any `json:"out_of_pocket"` //: [],
	Limitations []any `json:"limitations"`   //: []
}

func (s *InsuranceEligibilityService) GetFullReport(ctx context.Context, id int64) (*InsuranceEligibilityFullReport, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "get insurance eligibility full report", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.patient_insurance_id", id)))
	defer span.End()

	out := &InsuranceEligibilityFullReport{}

	res, err := s.client.request(ctx, http.MethodGet, fmt.Sprintf("/patient_insurances/%d/eligibility_full_report/", id), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

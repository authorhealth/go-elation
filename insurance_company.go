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

type InsuranceCompanyServicer interface {
	Create(ctx context.Context, create *InsuranceCompanyCreate) (*InsuranceCompany, *http.Response, error)
	Find(ctx context.Context, opts *FindInsuranceCompaniesOptions) (*Response[[]*InsuranceCompany], *http.Response, error)
	Get(ctx context.Context, id int64) (*InsuranceCompany, *http.Response, error)
	Update(ctx context.Context, id int64, update *InsuranceCompanyUpdate) (*InsuranceCompany, *http.Response, error)
	Delete(ctx context.Context, id int64) (*http.Response, error)
}

var _ InsuranceCompanyServicer = (*InsuranceCompanyService)(nil)

type InsuranceCompanyService struct {
	client *HttpClient
}

type InsuranceCompanyCreate struct {
	Practice           int64  `json:"practice"`
	Carrier            string `json:"carrier"`
	Address            string `json:"address,omitempty"`
	Suite              string `json:"suite,omitempty"`
	City               string `json:"city,omitempty"`
	State              string `json:"state,omitempty"`
	Zip                string `json:"zip,omitempty"`
	Phone              string `json:"phone,omitempty"`
	Extension          string `json:"extension,omitempty"`
	PayerID            string `json:"payer_id,omitempty"`
	ExternalVendorID   string `json:"external_vendor_id,omitempty"`
	IsCredentialed     *bool  `json:"is_credentialed,omitempty"`
	Aliases            string `json:"aliases,omitempty"`
	InsuranceType      string `json:"insurance_type,omitempty"`
	EligibilityPayerID string `json:"eligibility_payer_id,omitempty"`
}

func (s *InsuranceCompanyService) Create(ctx context.Context, create *InsuranceCompanyCreate) (*InsuranceCompany, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "create insurance company", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &InsuranceCompany{}

	res, err := s.client.request(ctx, http.MethodPost, "/insurance_companies", nil, create, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

type InsuranceCompany struct {
	ID                 int64   `json:"id"`
	Practice           int64   `json:"practice"`
	Carrier            string  `json:"carrier"`
	Address            string  `json:"address"`
	Suite              string  `json:"suite"`
	City               string  `json:"city"`
	State              string  `json:"state"`
	Zip                string  `json:"zip"`
	Phone              string  `json:"phone"`
	Extension          string  `json:"extension"`
	CreatedDate        string  `json:"created_date"`
	DeletedDate        string  `json:"deleted_date"`
	Patients           []int64 `json:"patients"`
	Plans              []int64 `json:"plans"`
	PayerID            string  `json:"payer_id"`
	ExternalVendorID   string  `json:"external_vendor_id"`
	IsCredentialed     bool    `json:"is_credentialed"`
	Aliases            string  `json:"aliases"`
	InsuranceType      string  `json:"insurance_type"`
	EligibilityPayerID string  `json:"eligibility_payer_id"`
}

type FindInsuranceCompaniesOptions struct {
	*Pagination

	Practice []int64 `url:"practice,omitempty"`
	Carrier  string  `url:"carrier,omitempty"`
}

func (s *InsuranceCompanyService) Find(ctx context.Context, opts *FindInsuranceCompaniesOptions) (*Response[[]*InsuranceCompany], *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "find insurance companies", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &Response[[]*InsuranceCompany]{}

	res, err := s.client.request(ctx, http.MethodGet, "/insurance_companies", opts, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *InsuranceCompanyService) Get(ctx context.Context, id int64) (*InsuranceCompany, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "get insurance company", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.insurance_company_id", id)))
	defer span.End()

	out := &InsuranceCompany{}

	res, err := s.client.request(ctx, http.MethodGet, "/insurance_companies/"+strconv.FormatInt(id, 10), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

type InsuranceCompanyUpdate struct {
	Carrier            string `json:"carrier"`
	Address            string `json:"address,omitempty"`
	Suite              string `json:"suite,omitempty"`
	City               string `json:"city,omitempty"`
	State              string `json:"state,omitempty"`
	Zip                string `json:"zip,omitempty"`
	Phone              string `json:"phone,omitempty"`
	Extension          string `json:"extension,omitempty"`
	PayerID            string `json:"payer_id,omitempty"`
	ExternalVendorID   string `json:"external_vendor_id,omitempty"`
	IsCredentialed     *bool  `json:"is_credentialed,omitempty"`
	Aliases            string `json:"aliases,omitempty"`
	InsuranceType      string `json:"insurance_type,omitempty"`
	EligibilityPayerID string `json:"eligibility_payer_id,omitempty"`
}

func (s *InsuranceCompanyService) Update(ctx context.Context, id int64, update *InsuranceCompanyUpdate) (*InsuranceCompany, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "update insurance company", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.insurance_company_id", id)))
	defer span.End()

	out := &InsuranceCompany{}

	res, err := s.client.request(ctx, http.MethodPatch, "/insurance_companies/"+strconv.FormatInt(id, 10), nil, update, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *InsuranceCompanyService) Delete(ctx context.Context, id int64) (*http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "delete insurance company", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.insurance_company_id", id)))
	defer span.End()

	res, err := s.client.request(ctx, http.MethodDelete, "/insurance_companies/"+strconv.FormatInt(id, 10), nil, nil, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return res, fmt.Errorf("making request: %w", err)
	}

	return res, nil
}

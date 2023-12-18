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

type InsurancePlanServicer interface {
	Create(ctx context.Context, create *InsurancePlanCreate) (*InsurancePlan, *http.Response, error)
	Find(ctx context.Context, opts *FindInsurancePlansOptions) (*Response[[]*InsurancePlan], *http.Response, error)
	Get(ctx context.Context, id int64) (*InsurancePlan, *http.Response, error)
	Update(ctx context.Context, id int64, update *InsurancePlanUpdate) (*InsurancePlan, *http.Response, error)
	Delete(ctx context.Context, id int64) (*http.Response, error)
}

var _ InsurancePlanServicer = (*InsurancePlanService)(nil)

type InsurancePlanService struct {
	client *Client
}

type InsurancePlanCreate struct {
	Practice         int64  `json:"practice"`
	InsuranceCompany int64  `json:"insurance_company"`
	Name             string `json:"name"`
}

func (s *InsurancePlanService) Create(ctx context.Context, create *InsurancePlanCreate) (*InsurancePlan, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "create insurance plan", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &InsurancePlan{}

	res, err := s.client.request(ctx, http.MethodPost, "/insurance_plans", nil, create, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

type InsurancePlan struct {
	ID               int64   `json:"id"`
	Practice         int64   `json:"practice"`
	InsuranceCompany int64   `json:"insurance_company"`
	Name             string  `json:"name"`
	CreatedDate      string  `json:"created_date"`
	DeletedDate      string  `json:"deleted_date"`
	Patients         []int64 `json:"patients"`
}

type FindInsurancePlansOptions struct {
	*Pagination

	Practice         []int64 `url:"practice,omitempty"`
	InsuranceCompany []int64 `url:"insurance_company,omitempty"`
}

func (s *InsurancePlanService) Find(ctx context.Context, opts *FindInsurancePlansOptions) (*Response[[]*InsurancePlan], *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "find insurance plans", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &Response[[]*InsurancePlan]{}

	res, err := s.client.request(ctx, http.MethodGet, "/insurance_plans", opts, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *InsurancePlanService) Get(ctx context.Context, id int64) (*InsurancePlan, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "get insurance plan", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.insurance_plan_id", id)))
	defer span.End()

	out := &InsurancePlan{}

	res, err := s.client.request(ctx, http.MethodGet, "/insurance_plans/"+strconv.FormatInt(id, 10), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

type InsurancePlanUpdate struct {
	Name string `json:"name"`
}

func (s *InsurancePlanService) Update(ctx context.Context, id int64, update *InsurancePlanUpdate) (*InsurancePlan, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "update insurance plan", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.insurance_plan_id", id)))
	defer span.End()

	out := &InsurancePlan{}

	res, err := s.client.request(ctx, http.MethodPatch, "/insurance_plans/"+strconv.FormatInt(id, 10), nil, update, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *InsurancePlanService) Delete(ctx context.Context, id int64) (*http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "delete insurance plan", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.insurance_plan_id", id)))
	defer span.End()

	res, err := s.client.request(ctx, http.MethodDelete, "/insurance_plans/"+strconv.FormatInt(id, 10), nil, nil, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return res, fmt.Errorf("making request: %w", err)
	}

	return res, nil
}

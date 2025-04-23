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

type AllergyDocumentationServicer interface {
	Find(ctx context.Context, opts *FindAllergiesDocumentationOptions) (*Response[[]*AllergyDocumentation], *http.Response, error)
	Get(ctx context.Context, id int64) (*AllergyDocumentation, *http.Response, error)
}

var _ AllergyDocumentationServicer = (*AllergyDocumentationService)(nil)

type AllergyDocumentationService struct {
	client *HTTPClient
}

type AllergyDocumentation struct {
	ID          int64      `json:"id"`
	Patient     int64      `json:"patient"`
	CreatedDate time.Time  `json:"created_date"`
	DeletedDate *time.Time `json:"deleted_date"`
}

type FindAllergiesDocumentationOptions struct {
	*Pagination

	Patient []int64 `url:"patient,omitempty"`
}

func (s *AllergyDocumentationService) Find(ctx context.Context, opts *FindAllergiesDocumentationOptions) (*Response[[]*AllergyDocumentation], *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "find allergy documentation", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &Response[[]*AllergyDocumentation]{}

	res, err := s.client.request(ctx, http.MethodGet, "/allergy_documentation", opts, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *AllergyDocumentationService) Get(ctx context.Context, id int64) (*AllergyDocumentation, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "get allergy documentation", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.allergy_documentation_id", id)))
	defer span.End()

	out := &AllergyDocumentation{}

	res, err := s.client.request(ctx, http.MethodGet, "/allergy_documentation/"+strconv.FormatInt(id, 10), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

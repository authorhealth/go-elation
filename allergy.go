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

type AllergyServicer interface {
	Find(ctx context.Context, opts *FindAllergiesOptions) (*Response[[]*Allergy], *http.Response, error)
	Get(ctx context.Context, id int64) (*Allergy, *http.Response, error)
}

var _ AllergyServicer = (*AllergyService)(nil)

type AllergyService struct {
	client *HTTPClient
}

type Allergy struct {
	ID          int64      `json:"id"`
	Status      string     `json:"status"`
	StartDate   time.Time  `json:"start_date"`
	Reaction    string     `json:"reaction"`
	Name        string     `json:"name"`
	Severity    string     `json:"severity"`
	MedispanID  *int64     `json:"medispanid"`
	MedispanNID *int64     `json:"medispandnid"`
	Patient     int64      `json:"patient"`
	CreatedDate time.Time  `json:"created_date"`
	DeletedDate *time.Time `json:"deleted_date"`
}

type FindAllergiesOptions struct {
	*Pagination

	Patient []int64 `url:"patient,omitempty"`
}

func (s *AllergyService) Find(ctx context.Context, opts *FindAllergiesOptions) (*Response[[]*Allergy], *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "find allergies", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &Response[[]*Allergy]{}

	res, err := s.client.request(ctx, http.MethodGet, "/allergies", opts, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *AllergyService) Get(ctx context.Context, id int64) (*Allergy, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "get allergy", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.allergy_id", id)))
	defer span.End()

	out := &Allergy{}

	res, err := s.client.request(ctx, http.MethodGet, "/allergies/"+strconv.FormatInt(id, 10), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

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

type PracticeServicer interface {
	Find(ctx context.Context, opts *FindPracticesOptions) (*Response[[]*Practice], *http.Response, error)
	Get(ctx context.Context, id int64) (*Practice, *http.Response, error)
}

var _ PracticeServicer = (*PracticeService)(nil)

type PracticeService struct {
	client *HTTPClient
}

type Practice struct {
	ID               int64                      `json:"id"`
	Name             string                     `json:"name"`
	AddressLine1     string                     `json:"address_line1"`
	AddressLine2     string                     `json:"address_line2"`
	City             string                     `json:"city"`
	State            string                     `json:"state"`
	Zip              string                     `json:"zip"`
	Timezone         string                     `json:"timezone"`
	ElationRootOid   string                     `json:"elation_root_oid"`
	Employers        []*PracticeEmployer        `json:"employers"`
	Physicians       []int64                    `json:"physicians"`
	ServiceLocations []*PracticeServiceLocation `json:"service_locations"`
	Metadata         any                        `json:"metadata"`
	Status           string                     `json:"status"`
}

type PracticeEmployer struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PracticeServiceLocation struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	IsPrimary      bool      `json:"is_primary"`
	PlaceOfService string    `json:"place_of_service"`
	AddressLine1   string    `json:"address_line1"`
	AddressLine2   string    `json:"address_line2"`
	City           string    `json:"city"`
	State          string    `json:"state"`
	Zip            string    `json:"zip"`
	Phone          string    `json:"phone"`
	Email          any       `json:"email"`
	Fax            string    `json:"fax"`
	Practice       int64     `json:"practice"`
	CreatedDate    time.Time `json:"created_date"`
	DeletedDate    any       `json:"deleted_date"`
	Status         string    `json:"status"`
}

type FindPracticesOptions struct {
	*Pagination
}

func (s *PracticeService) Find(ctx context.Context, opts *FindPracticesOptions) (*Response[[]*Practice], *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "find practices", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &Response[[]*Practice]{}

	res, err := s.client.request(ctx, http.MethodGet, "/practices", opts, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *PracticeService) Get(ctx context.Context, id int64) (*Practice, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "get practice", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.Practice_id", id)))
	defer span.End()

	out := &Practice{}

	res, err := s.client.request(ctx, http.MethodGet, "/practices/"+strconv.FormatInt(id, 10), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

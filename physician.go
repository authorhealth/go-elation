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

type PhysicianServicer interface {
	Find(ctx context.Context, opts *FindPhysiciansOptions) (*Response[[]*Physician], *http.Response, error)
	Get(ctx context.Context, id int64) (*Physician, *http.Response, error)
}

var _ PhysicianServicer = (*PhysicianService)(nil)

type PhysicianService struct {
	client *Client
}

type Physician struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Npi          string `json:"npi"`
	License      string `json:"license"`
	LicenseState string `json:"license_state"`
	Credentials  string `json:"credentials"`
	Specialty    string `json:"specialty"`
	Email        string `json:"email"`
	UserID       int    `json:"user_id"`
	Practice     int    `json:"practice"`
	IsActive     bool   `json:"is_active"`
	Metadata     any    `json:"metadata"`
}

type FindPhysiciansOptions struct {
	*Pagination

	FirstName string `url:"first_name,omitempty"`
	LastName  string `url:"last_name,omitempty"`
	NPI       string `url:"npi,omitempty"`
}

func (s *PhysicianService) Find(ctx context.Context, opts *FindPhysiciansOptions) (*Response[[]*Physician], *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "find physicians", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &Response[[]*Physician]{}

	res, err := s.client.request(ctx, http.MethodGet, "/physicians", opts, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *PhysicianService) Get(ctx context.Context, id int64) (*Physician, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "get physician", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.physician_id", id)))
	defer span.End()

	out := &Physician{}

	res, err := s.client.request(ctx, http.MethodGet, "/physicians/"+strconv.FormatInt(id, 10), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

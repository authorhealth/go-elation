package elation

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type ServiceLocationServicer interface {
	Find(ctx context.Context, opts *FindServiceLocationOptions) (*Response[[]*ServiceLocation], *http.Response, error)
}

var _ ServiceLocationServicer = (*ServiceLocationService)(nil)

type ServiceLocationService struct {
	client *HttpClient
}

type ServiceLocation struct {
	AddressLine1   string     `json:"address_line1"`
	AddressLine2   string     `json:"address_line2"`
	City           string     `json:"city"`
	CreatedDate    time.Time  `json:"created_date"`
	DeletedDate    *time.Time `json:"deleted_date"`
	Email          string     `json:"email"`
	Fax            string     `json:"fax"`
	ID             int64      `json:"id"`
	IsPrimary      bool       `json:"is_primary"`
	Name           string     `json:"name"`
	Phone          string     `json:"phone"`
	PlaceOfService string     `json:"place_of_service"`
	Practice       int64      `json:"practice"`
	State          string     `json:"state"`
	Status         string     `json:"status"`
	Zip            string     `json:"zip"`
}

type FindServiceLocationOptions struct {
	*Pagination
}

func (s *ServiceLocationService) Find(ctx context.Context, opts *FindServiceLocationOptions) (*Response[[]*ServiceLocation], *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "find service locations", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &Response[[]*ServiceLocation]{}

	res, err := s.client.request(ctx, http.MethodGet, "/service_locations", opts, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

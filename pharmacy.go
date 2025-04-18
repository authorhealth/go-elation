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

type PharmacyServicer interface {
	Get(ctx context.Context, ncpdpid string) (*Pharmacy, *http.Response, error)
}

var _ PharmacyServicer = (*PharmacyService)(nil)

type PharmacyService struct {
	client *HTTPClient
}

type Pharmacy struct {
	ID              int64     `json:"id"`
	NCPDPID         string    `json:"ncpdpid"`
	StoreName       string    `json:"store_name"`
	AddressLine1    string    `json:"address_line1"`
	AddressLine2    string    `json:"address_line2"`
	City            string    `json:"city"`
	State           string    `json:"state"`
	Zip             string    `json:"zip"`
	PhonePrimary    string    `json:"phone_primary"`
	Fax             string    `json:"fax"`
	NPI             string    `json:"npi"`
	ActiveStartTime time.Time `json:"active_start_time"`
	ActiveEndTime   time.Time `json:"active_end_time"`
	SpecialityTypes string    `json:"specialty_types"`
}

func (s *PharmacyService) Get(ctx context.Context, ncpdpid string) (*Pharmacy, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "get pharmacy", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.String("elation.pharmacy_ncpdpid", ncpdpid)))
	defer span.End()

	out := &Pharmacy{}

	res, err := s.client.request(ctx, http.MethodGet, fmt.Sprintf("/pharmacies/%s", ncpdpid), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

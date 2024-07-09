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

type ContactServicer interface {
	Get(ctx context.Context, id int64) (*Contact, *http.Response, error)
	List(ctx context.Context, opts *ListContactsOptions) (*Response[[]*Contact], *http.Response, error)
}

var _ ContactServicer = (*ContactService)(nil)

type ContactService struct {
	client *HTTPClient
}

type Contact struct {
	ID                   int64          `json:"id"`
	AcceptingInsurance   bool           `json:"accepting_insurance"`
	AcceptingNewPatients bool           `json:"accepting_new_patients"`
	Address              string         `json:"address"`
	BackOfficePhone      string         `json:"back_office_phone"`
	CellPhone            string         `json:"cell_phone"`
	City                 string         `json:"city"`
	ContactType          string         `json:"contact_type"`
	CreatedDate          time.Time      `json:"created_date"`
	Credentials          string         `json:"credentials"`
	DeletedDate          *time.Time     `json:"deleted_date"`
	DirectAddress        string         `json:"direct_address"`
	Email                string         `json:"email"`
	Fax                  string         `json:"fax"`
	FirstName            string         `json:"first_name"`
	IsElationConfirmed   bool           `json:"is_elation_confirmed"`
	IsVerified           bool           `json:"is_verified"`
	LastName             string         `json:"last_name"`
	MiddleName           string         `json:"middle_name"`
	NPI                  string         `json:"npi"`
	OrgName              string         `json:"org_name"`
	OtherSpecialties     []any          `json:"other_specialties"`
	Phone                string         `json:"phone"`
	Practice             int64          `json:"practice"`
	Specialty            map[string]any `json:"specialty"`
	State                string         `json:"state"`
	Suite                string         `json:"suite"`
	User                 int64          `json:"user"`
	Zip                  string         `json:"zip"`
}

type ListContactsOptions struct {
	*Pagination

	NPI string `url:"npi,omitempty"`
}

func (s *ContactService) List(ctx context.Context, opts *ListContactsOptions) (*Response[[]*Contact], *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "find contacts", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &Response[[]*Contact]{}

	res, err := s.client.request(ctx, http.MethodGet, "/contacts", opts, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *ContactService) Get(ctx context.Context, id int64) (*Contact, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "get contact", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.contact_id", id)))
	defer span.End()

	out := &Contact{}

	res, err := s.client.request(ctx, http.MethodGet, "/contacts/"+strconv.FormatInt(id, 10), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

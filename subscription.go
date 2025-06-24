package elation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type SubscriptionServicer interface {
	Find(ctx context.Context) ([]*Subscription, *http.Response, error)
	Subscribe(ctx context.Context, opts *Subscribe) (*Subscription, *http.Response, error)
	Delete(ctx context.Context, id int64) (*http.Response, error)
}

var _ SubscriptionServicer = (*SubscriptionService)(nil)

type SubscriptionService struct {
	client *HTTPClient
}

type Subscription struct {
	ID            int64                 `json:"id"`
	Resource      Resource              `json:"resource"`
	Target        string                `json:"target"`
	CreatedDate   SubscriptionJSONDate  `json:"created_date"`
	DeletedDate   *SubscriptionJSONDate `json:"deleted_date"`
	SigningPubKey string                `json:"signing_pub_key"`
}

type Resource string

const (
	ResourceAppointments            Resource = "appointments"
	ResourceAllergies               Resource = "allergies"
	ResourceAllergyDocumentation    Resource = "allergy_documentation"
	ResourceDiscontinuedMedications Resource = "discontinued_medications"
	ResourceLetters                 Resource = "letters"
	ResourceMedications             Resource = "medications"
	ResourcePatients                Resource = "patients"
	ResourcePhysicians              Resource = "physicians"
	ResourceProblems                Resource = "problems"
)

func (r Resource) String() string {
	return string(r)
}

type SubscriptionJSONDate time.Time

func (s *SubscriptionJSONDate) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Time(*s).Format("2006-01-02T15:04:05") + "\""), nil
}

func (s *SubscriptionJSONDate) UnmarshalJSON(b []byte) error {
	t, err := time.Parse("2006-01-02T15:04:05", strings.Trim(string(b), "\""))
	if err != nil {
		return err
	}

	*s = SubscriptionJSONDate(t)

	return nil
}

func (s *SubscriptionService) Find(ctx context.Context) ([]*Subscription, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "find subscriptions", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &Response[[]*Subscription]{}

	// The trailing slash in the path is required.
	res, err := s.client.request(ctx, http.MethodGet, "/app/subscriptions/", nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out.Results, res, nil
}

type Subscribe struct {
	Resource   Resource        `json:"resource"`
	Target     string          `json:"target"`
	Properties json.RawMessage `json:"properties"`
}

func (s *SubscriptionService) Subscribe(ctx context.Context, subscribe *Subscribe) (*Subscription, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "subscribe", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &Subscription{}

	// The trailing slash in the path is required.
	res, err := s.client.request(ctx, http.MethodPost, "/app/subscriptions/", nil, subscribe, out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

type DeleteSubscriptionOptions struct {
	ID int64 `url:"id"`
}

func (s *SubscriptionService) Delete(ctx context.Context, id int64) (*http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "delete subscription", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.subscription_id", id)))
	defer span.End()

	// The trailing slash in the path is required.
	res, err := s.client.request(ctx, http.MethodDelete, "/app/subscriptions/"+strconv.FormatInt(id, 10)+"/", nil, nil, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return res, fmt.Errorf("making request: %w", err)
	}

	return res, nil
}

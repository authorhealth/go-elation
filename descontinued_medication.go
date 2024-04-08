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

type DiscontinuedMedicationServicer interface {
	Create(ctx context.Context, create *DiscontinuedMedicationCreate) (*DiscontinuedMedication, *http.Response, error)
	Find(ctx context.Context, opts *FindDiscontinuedMedicationsOptions) (*Response[[]*DiscontinuedMedication], *http.Response, error)
	Get(ctx context.Context, id int64) (*DiscontinuedMedication, *http.Response, error)
	Update(ctx context.Context, id int64, update *DiscontinuedMedicationUpdate) (*DiscontinuedMedication, *http.Response, error)
}

var _ DiscontinuedMedicationServicer = (*DiscontinuedMedicationService)(nil)

type DiscontinuedMedicationService struct {
	client *Client
}

type DiscontinuedMedicationCreate struct {
	MedOrder             int64  `json:"med_order"`
	DiscontinueDate      string `json:"discontinue_date,omitempty"`
	Reason               string `json:"reason,omitempty"`
	IsDocumented         bool   `json:"is_documented"`
	DocumentingPersonnel int64  `json:"documenting_personnel,omitempty"`
}

func (s *DiscontinuedMedicationService) Create(ctx context.Context, create *DiscontinuedMedicationCreate) (*DiscontinuedMedication, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "create discontinued medication", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &DiscontinuedMedication{}

	res, err := s.client.request(ctx, http.MethodPost, "/discontinued_medications", nil, create, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

type DiscontinuedMedication struct {
	ID                   int64                             `json:"id"`
	LastMedicationOrder  int64                             `json:"last_medication_order"`
	MedOrder             int64                             `json:"med_order"`
	Thread               int64                             `json:"thread"`
	DiscontinueDate      string                            `json:"discontinue_date"`
	Reason               string                            `json:"reason"`
	PrescribingPhysician int64                             `json:"prescribing_physician"`
	SupervisingPhysician int64                             `json:"supervising_physician"`
	IsDocumented         bool                              `json:"is_documented"`
	Medication           *DiscontinuedMedicationMedication `json:"medication"`
	DocumentingPersonnel int64                             `json:"documenting_personnel"`
	DocumentDate         time.Time                         `json:"document_date"`
	ChartDate            time.Time                         `json:"chart_date"`
	CreatedDate          time.Time                         `json:"created_date"`
	DeletedDate          *time.Time                        `json:"deleted_date"`
	SignedDate           *time.Time                        `json:"signed_date"`
	SignedBy             int64                             `json:"signed_by"`
	Patient              int64                             `json:"patient"`
	Practice             int64                             `json:"practice"`
}

type DiscontinuedMedicationMedication struct {
	ID             int64    `json:"id"`
	RxnormCuis     []string `json:"rxnorm_cuis"`
	NDCs           []string `json:"ndcs"`
	Name           string   `json:"name"`
	BrandName      string   `json:"brand_name"`
	GenericName    string   `json:"generic_name"`
	IsControlled   bool     `json:"is_controlled"`
	MedicationType string   `json:"type"`
	Route          string   `json:"route"`
	Strength       string   `json:"strength"`
	Form           string   `json:"form"`
}

type FindDiscontinuedMedicationsOptions struct {
	*Pagination

	Patient         []int64   `url:"patient,omitempty"`
	Practice        []int64   `url:"practice,omitempty"`
	DocumentDateLT  time.Time `url:"document_date__lt,omitempty"`
	DocumentDateGT  time.Time `url:"document_date__gt,omitempty"`
	DocumentDateLTE time.Time `url:"document_date__lte,omitempty"`
	DocumentDateGTE time.Time `url:"document_date__gte,omitempty"`
}

func (s *DiscontinuedMedicationService) Find(ctx context.Context, opts *FindDiscontinuedMedicationsOptions) (*Response[[]*DiscontinuedMedication], *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "find discontinued medications", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &Response[[]*DiscontinuedMedication]{}

	res, err := s.client.request(ctx, http.MethodGet, "/discontinued_medications", opts, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *DiscontinuedMedicationService) Get(ctx context.Context, id int64) (*DiscontinuedMedication, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "get discontinued medication", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.discontinued_medication_id", id)))
	defer span.End()

	out := &DiscontinuedMedication{}

	res, err := s.client.request(ctx, http.MethodGet, "/discontinued_medications/"+strconv.FormatInt(id, 10), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

type DiscontinuedMedicationUpdate struct {
	Reason               string `json:"reason,omitempty"`
	DiscontinueDate      string `json:"discontinue_date,omitempty"`
	DocumentingPersonnel int64  `json:"documenting_personnel,omitempty"`
}

func (s *DiscontinuedMedicationService) Update(ctx context.Context, id int64, update *DiscontinuedMedicationUpdate) (*DiscontinuedMedication, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "update discontinued medication", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.discontinued_medication_id", id)))
	defer span.End()

	out := &DiscontinuedMedication{}

	res, err := s.client.request(ctx, http.MethodPatch, "/discontinued_medications/"+strconv.FormatInt(id, 10), nil, update, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

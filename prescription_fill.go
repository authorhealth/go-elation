package elation

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/civil"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type PrescriptionFillServicer interface {
	Find(ctx context.Context, opts *FindPrescriptionFillsOptions) (*Response[[]*PrescriptionFill], *http.Response, error)
	Get(ctx context.Context, id int64) (*PrescriptionFill, *http.Response, error)
}

var _ PrescriptionFillServicer = (*PrescriptionFillService)(nil)

type PrescriptionFillService struct {
	client *HttpClient
}

type PrescriptionFill struct {
	ID               int64       `json:"id"`
	MedicationOrder  int64       `json:"medication_order"`
	PharmacyNCPDPID  string      `json:"pharmacy_ncpdpid"`
	FillStatus       string      `json:"fill_status"`
	FillDate         *civil.Date `json:"fill_date"`
	NoteFromPharmacy string      `json:"note_from_pharmacy"`
	Patient          int64       `json:"patient"`
	Practice         int64       `json:"practice"`
}

type FindPrescriptionFillsOptions struct {
	*Pagination

	Practice       []int64   `url:"practice,omitempty"`
	Patient        []int64   `url:"patient,omitempty"`
	FillDateIsNull *bool     `url:"fill_date__isnull,omitempty"`
	FillDateGTE    time.Time `url:"fill_date__gte,omitempty"`
	FillDateLTE    time.Time `url:"fill_date__lte,omitempty"`
	FillStatus     string    `url:"fill_status,omitempty"`
}

func (s *PrescriptionFillService) Find(ctx context.Context, opts *FindPrescriptionFillsOptions) (*Response[[]*PrescriptionFill], *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "find prescription fills", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &Response[[]*PrescriptionFill]{}

	res, err := s.client.request(ctx, http.MethodGet, "/prescription_fills", opts, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *PrescriptionFillService) Get(ctx context.Context, id int64) (*PrescriptionFill, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "get prescription fill", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.prescription_fill_id", id)))
	defer span.End()

	out := &PrescriptionFill{}

	res, err := s.client.request(ctx, http.MethodGet, "/prescription_fills/"+strconv.FormatInt(id, 10), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

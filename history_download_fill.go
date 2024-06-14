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

type HistoryDownloadFillServicer interface {
	Find(ctx context.Context, opts *FindHistoryDownloadFillsOptions) (*Response[[]*HistoryDownloadFill], *http.Response, error)
	Get(ctx context.Context, id int64) (*HistoryDownloadFill, *http.Response, error)
}

var _ HistoryDownloadFillServicer = (*HistoryDownloadFillService)(nil)

type HistoryDownloadFillService struct {
	client *Client
}

type HistoryDownloadFill struct {
	ID                        int64      `json:"id"`
	Patient                   int64      `json:"patient"`
	Practice                  int64      `json:"practice"`
	CreateDate                time.Time  `json:"create_date"`
	MedicationOrder           int64      `json:"medication_order"`
	HistoryDownload           int64      `json:"history_download"`
	Medication                int64      `json:"medication"`
	MedicationDescription     string     `json:"medication_description"`
	Quantity                  string     `json:"quantity"`
	QuantityUnit              string     `json:"quantity_unit"`
	QuantityNote              string     `json:"quantity_note"`
	DaysSupply                string     `json:"days_supply"`
	Directions                string     `json:"directions"`
	Note                      string     `json:"note"`
	Substitutions             string     `json:"substitutions"`
	LastFillDate              time.Time  `json:"last_fill_date"`
	WrittenDate               *time.Time `json:"written_date"`
	Refills                   []int64    `json:"refills"`
	Diagnoses                 []int64    `json:"diagnoses"`
	PriorAuth                 string     `json:"prior_auth"`
	PriorAuthQualifier        string     `json:"prior_auth_qualifier"`
	PharmacyNCPDPID           string     `json:"pharmacy_ncpdpid"`
	Prescriber                int64      `json:"prescriber"`
	PrescribingPhysician      int64      `json:"prescribing_physician"`
	RequestingPhysicianAction string     `json:"requesting_physician_action"`
	IsActive                  bool       `json:"is_active"`
	RequestingPhysician       int64      `json:"requesting_physician"`
}

type FindHistoryDownloadFillsOptions struct {
	*Pagination

	Practice           []int64   `url:"practice,omitempty"`
	Patient            []int64   `url:"patient,omitempty"`
	HistoryDownload    string    `url:"history_download,omitempty"`
	LastFillDateGTE    time.Time `url:"last_fill_date__gte,omitempty"`
	LastFillDateLTE    time.Time `url:"last_fill_date__lte,omitempty"`
	WrittenDateLTE     time.Time `url:"written_date__lte,omitempty"`
	WrittenDateGTE     time.Time `url:"written_date__gte,omitempty"`
	LastFillDateIsNull *bool     `url:"last_fill_date__isnull,omitempty"`
	WrittenDateIsNull  *bool     `url:"written_date__isnull,omitempty"`
}

func (s *HistoryDownloadFillService) Find(ctx context.Context, opts *FindHistoryDownloadFillsOptions) (*Response[[]*HistoryDownloadFill], *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "find history download fills", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &Response[[]*HistoryDownloadFill]{}

	res, err := s.client.request(ctx, http.MethodGet, "/medication_history_download_fills", opts, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *HistoryDownloadFillService) Get(ctx context.Context, id int64) (*HistoryDownloadFill, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "get history download fill", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.medication_history_download_fill_id", id)))
	defer span.End()

	out := &HistoryDownloadFill{}

	res, err := s.client.request(ctx, http.MethodGet, "/medication_history_download_fills/"+strconv.FormatInt(id, 10), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

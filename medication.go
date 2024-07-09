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

type MedicationServicer interface {
	Create(ctx context.Context, create *PatientMedicationCreate) (*PatientMedication, *http.Response, error)
	Find(ctx context.Context, opts *FindPatientMedicationsOptions) (*Response[[]*PatientMedication], *http.Response, error)
	Get(ctx context.Context, id int64) (*PatientMedication, *http.Response, error)
}

var _ MedicationServicer = (*MedicationService)(nil)

type MedicationService struct {
	client *HTTPClient
}

type PatientMedication struct {
	ID                   int64                         `json:"id"`
	Patient              int64                         `json:"patient"`
	Medication           *Medication                   `json:"medication"`
	Fulfillment          *PatientMedicationFulfillment `json:"fulfillment"`
	Thread               *PatientMedicationThread      `json:"thread"`
	OrderType            string                        `json:"order_type"`
	Qty                  string                        `json:"qty"`
	QtyUnits             string                        `json:"qty_units"`
	AuthRefills          int                           `json:"auth_refills"`
	NumSamples           any                           `json:"num_samples"`
	PrescribingPhysician int                           `json:"prescribing_physician"`
	Practice             int64                         `json:"practice"`
	MedicationType       string                        `json:"medication_type"`
	Notes                string                        `json:"notes"`
	StartDate            string                        `json:"start_date"`
	IsDocMed             bool                          `json:"is_doc_med"`
	DocumentingPersonnel int                           `json:"documenting_personnel"`
	LastModified         string                        `json:"last_modified"`
	PharmacyInstructions string                        `json:"pharmacy_instructions"`
	Directions           string                        `json:"directions"`
	DocumentDate         time.Time                     `json:"document_date"`
	ChartDate            time.Time                     `json:"chart_date"`
	ShowInChartFeed      bool                          `json:"show_in_chart_feed"`
	SignedDate           time.Time                     `json:"signed_date"`
	SignedBy             int                           `json:"signed_by"`
	CreatedDate          time.Time                     `json:"created_date"`
	DeletedDate          *time.Time                    `json:"deleted_date"`
	Icd10Codes           []*PatientMedicationICD10Code `json:"icd10_codes"`
}

type PatientMedicationCreate struct {
	AuthRefills          int                                 `json:"auth_refills"`
	IsDocMed             bool                                `json:"is_doc_med"`
	Medication           *PatientMedicationCreateMedication  `json:"medication"`
	MedicationType       string                              `json:"medication_type"`
	OrderType            string                              `json:"order_type"`
	Patient              int64                               `json:"patient"`
	Practice             int64                               `json:"practice"`
	PrescribingPhysician int                                 `json:"prescribing_physician"`
	Qty                  string                              `json:"qty"`
	QtyUnits             string                              `json:"qty_units"`
	StartDate            string                              `json:"start_date"`
	Icd10Codes           []*PatientMedicationCreateICD10Code `json:"icd10_codes"`
	Thread               *PatientMedicationCreateThread      `json:"thread"`
	Directions           string                              `json:"directions,omitempty"`
	Notes                string                              `json:"notes,omitempty"`
	PharmacyInstructions string                              `json:"pharmacy_instructions,omitempty"`
	NumSamples           string                              `json:"num_samples,omitempty"`
	DocumentDate         *time.Time                          `json:"document_date,omitempty"`
	DocumentingPersonnel int                                 `json:"documenting_personnel,omitempty"`
}

type PatientMedicationCreateMedication struct {
	ID         int64    `json:"id,omitempty"`
	NDCs       []string `json:"ndcs,omitempty"`
	RxnormCuis []string `json:"rxnorm_cuis,omitempty"`
}

type PatientMedicationCreateICD10Code struct {
	Code string `json:"code"`
}

type PatientMedicationCreateThread struct {
	ID          int  `json:"id,omitempty"`
	IsPermanent bool `json:"is_permanent"`
}

type Medication struct {
	ID           int64    `json:"id"`
	IsControlled bool     `json:"is_controlled"`
	Strength     string   `json:"strength"`
	Name         string   `json:"name"`
	Form         string   `json:"form"`
	Route        string   `json:"route"`
	GenericName  string   `json:"generic_name"`
	BrandName    string   `json:"brand_name"`
	Type         string   `json:"type"`
	RxnormCuis   []string `json:"rxnorm_cuis"`
}

type PatientMedicationFulfillment struct {
	Detail          any                                          `json:"detail"`
	PharmacyNcpdpid string                                       `json:"pharmacy_ncpdpid"`
	ServiceLocation *PatientMedicationFulfillmentServiceLocation `json:"service_location"`
	State           string                                       `json:"state"`
	TimeCompleted   time.Time                                    `json:"time_completed"`
	Type            string                                       `json:"type"`
}

type PatientMedicationFulfillmentServiceLocation struct {
	ID             int64      `json:"id"`
	Name           string     `json:"name"`
	IsPrimary      bool       `json:"is_primary"`
	PlaceOfService string     `json:"place_of_service"`
	AddressLine1   string     `json:"address_line1"`
	AddressLine2   string     `json:"address_line2"`
	City           string     `json:"city"`
	State          string     `json:"state"`
	Zip            string     `json:"zip"`
	Phone          string     `json:"phone"`
	CreatedDate    time.Time  `json:"created_date"`
	DeletedDate    *time.Time `json:"deleted_date"`
}

type PatientMedicationThread struct {
	ID          int    `json:"id"`
	DcDate      string `json:"dc_date"`
	IsPermanent bool   `json:"is_permanent"`
}

type PatientMedicationICD10Code struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type FindPatientMedicationsOptions struct {
	*Pagination

	Patient         int64     `url:"patient,omitempty"`
	Thread          int64     `url:"thread,omitempty"`
	Practice        int64     `url:"practice,omitempty"`
	Discontinued    bool      `url:"discontinued,omitempty"`
	Permanent       bool      `url:"permanent,omitempty"`
	DocumentDateLT  time.Time `url:"document_date__lt,omitempty"`
	DocumentDateGT  time.Time `url:"document_date__gt,omitempty"`
	DocumentDateLTE time.Time `url:"document_date__lte,omitempty"`
	DocumentDateGTE time.Time `url:"document_date__gte,omitempty"`
}

func (s *MedicationService) Create(ctx context.Context, create *PatientMedicationCreate) (*PatientMedication, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "create medication", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &PatientMedication{}

	res, err := s.client.request(ctx, http.MethodPost, "/medications", nil, create, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *MedicationService) Find(ctx context.Context, opts *FindPatientMedicationsOptions) (*Response[[]*PatientMedication], *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "find medications", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &Response[[]*PatientMedication]{}

	res, err := s.client.request(ctx, http.MethodGet, "/medications", opts, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *MedicationService) Get(ctx context.Context, id int64) (*PatientMedication, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "get medication", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.medication_id", id)))
	defer span.End()

	out := &PatientMedication{}

	res, err := s.client.request(ctx, http.MethodGet, "/medications/"+strconv.FormatInt(id, 10), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

package elation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const billExistError = "The visit note provided already has a bill associated with it."

var ErrBillExist = errors.New("bill already exists for visit note")

type BillServicer interface {
	Create(ctx context.Context, create *BillCreate) (*Bill, *http.Response, error)
}

var _ BillServicer = (*BillService)(nil)

type BillService struct {
	client *HTTPClient
}

type Bill struct {
	ID                   int64           `json:"id"`                      //: 65099661468,
	RefNumber            *string         `json:"ref_number"`              //: null,                        // string(50). required for PATCH that marks bill as processed.
	ServiceDate          time.Time       `json:"service_date"`            //: "2016-10-12T12:00:00Z",
	BillingDate          *time.Time      `json:"billing_date"`            //: null,                        // datetime(iso8601). required for PATCH that marks bill as processed.
	BillingStatus        string          `json:"billing_status"`          //: "Unbilled",
	BillingError         *string         `json:"billing_error"`           //: null,                        // string(200). required for PATCH that marks bill as failed.
	BillingRawError      *string         `json:"billing_raw_error"`       //: null,                        // longtext. optional for PATCH that marks bill as failed.
	Notes                string          `json:"notes"`                   //: "patient has not paid yet",
	CPTs                 []*BillCPT      `json:"cpts"`                    //: [{}],
	Payment              BillPayment     `json:"payment"`                 //: {"amount": "10.00","when_collected": "2016-10-12T22:11:01Z"}
	VisitNote            int64           `json:"visit_note"`              //: 64409108504,
	VisitNoteSignedDate  time.Time       `json:"visit_note_signed_date"`  //: "2016-10-12T22:11:01Z",
	VisitNoteDeletedDate *time.Time      `json:"visit_note_deleted_date"` //: null,
	ReferringProvider    *BillProvider   `json:"referring_provider"`      //: {},
	BillingProvider      *int64          `json:"billing_provider"`        //: 42120898,
	RenderingProvider    *int64          `json:"rendering_provider"`      //: 68382673,
	SupervisingProvider  *int64          `json:"supervising_provider"`    //: 52893234,
	OrderingProvider     *BillProvider   `json:"ordering_provider"`       //: {}
	ServiceLocation      ServiceLocation `json:"service_location"`        //: 141103949480183,
	Physician            int64           `json:"physician"`               //: 64811630594,
	Practice             int64           `json:"practice"`                //: 65540,
	Patient              int64           `json:"patient"`                 //: 64901939201,
	PriorAuthorization   *string         `json:"prior_authorization"`     //: "1234-ABC",
	Metadata             any             `json:"metadata"`                //: null,
	CreatedDate          time.Time       `json:"created_date"`            //: "2016-05-23T17:50:50Z",
	LastModifiedDate     time.Time       `json:"last_modified_date"`      //: "2016-10-12T22:39:46Z"
}

type BillCreate struct {
	ServiceLocation     int64            `json:"service_location"`               //: 10           // required
	VisitNote           int64            `json:"visit_note"`                     //: 64409108504, // required
	Patient             int64            `json:"patient"`                        //: 64901939201, // required
	Practice            int64            `json:"practice"`                       //: 65540, 		   // required
	Physician           int64            `json:"physician"`                      //: 64811630594, // required
	CPTs                []*BillCPTCreate `json:"cpts"`                           //: [{}],        // required
	BillingProvider     int64            `json:"billing_provider,omitempty"`     //: 42120898,
	RenderingProvider   int64            `json:"rendering_provider,omitempty"`   //: 68382673,
	SupervisingProvider int64            `json:"supervising_provider,omitempty"` //: 52893234,
	ReferringProvider   *BillProvider    `json:"referring_provider,omitempty"`   //: {},
	OrderingProvider    *BillProvider    `json:"ordering_provider,omitempty"`    //: {},
	PriorAuthorization  string           `json:"prior_authorization,omitempty"`  //: "1234-ABC",
	PaymentAmount       float64          `json:"payment_amount,omitempty"`       //: 10.00,
	Notes               string           `json:"notes,omitempty"`                //: "patient has not paid yet",
}

type BillDX struct {
	ICD10Code string `json:"icd10_code"`
}

type BillCPTCreate struct {
	CPT        string   `json:"cpt"`                 //: "99213",
	Modifiers  []string `json:"modifiers,omitempty"` //: ["10"],
	DXs        []BillDX `json:"dxs"`                 //: ["D23.4"],
	AltDXs     []string `json:"alt_dxs,omitempty"`   //: ["216.4"],
	UnitCharge string   `json:"unit_charge"`         //: "10.0",
	Units      string   `json:"units"`               //: "1.0"
}

type BillCPT struct {
	CPT        string   `json:"cpt"`                   //: "99213",
	Modifiers  []string `json:"modifiers,omitempty"`   //: ["10"],
	DXs        []string `json:"dxs"`                   //: ["D23.4"],
	AltDXs     []string `json:"alt_dxs,omitempty"`     //: ["216.4"],
	NDC        string   `json:"ndc,omitempty"`         // : "60575-4112-01", 11 digit NDC code, zero padded to 5-4-2 format
	NDCDose    string   `json:"ndc_dose,omitempty"`    // : "1.000", NDC units administered.  Precision to 3 decimal places
	NDCMeasure string   `json:"ndc_measure,omitempty"` // : "ML", NDC unit of measure: null, "F2", "GR", "ME", "ML" or "UN"
	UnitCharge string   `json:"unit_charge"`           //: "10.0", charge per CPT unit
	Units      string   `json:"units"`                 //: "1.0", CPT units
}

type BillPayment struct {
	Amount        string    `json:"amount"`         //: "10.00",
	WhenCollected time.Time `json:"when_collected"` //: "2016-10-12T22:11:01Z"
}

type BillProvider struct {
	Name  string `json:"name"`  //: "Beverly Crusher, MD (555-555-5555)",
	State string `json:"state"` //: "CA",
	NPI   string `json:"npi"`   //: "1701170117"
}

type BillServiceLocation struct {
	ID                 int64      `json:"id"`                    //: 13631735,
	Name               string     `json:"name"`                  //: "Elation North",
	IsPrimary          bool       `json:"is_primary"`            //: true,
	PlaceOfService     string     `json:"place_of_service"`      //: "Office",
	PlaceOfServiceCode string     `json:"place_of_service_code"` //: "11",
	AddressLine1       string     `json:"address_line1"`         //: "1234 First Practice Way",
	AddressLine2       string     `json:"address_line2"`         //: "",
	City               string     `json:"city"`                  //: "San Francisco",
	State              string     `json:"state"`                 //: "CA",
	Zip                string     `json:"zip"`                   //: "94114",
	Phone              string     `json:"phone"`                 //: "",
	CreatedDate        time.Time  `json:"created_date"`          //: "2017-08-28T22:46:14.445876Z",
	DeletedDate        *time.Time `json:"deleted_date"`          //: null
}

func (b *BillService) Create(ctx context.Context, create *BillCreate) (*Bill, *http.Response, error) {
	ctx, span := b.client.tracer.Start(ctx, "create bill", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	bill := &Bill{}

	res, err := b.client.request(ctx, http.MethodPost, "/bills", nil, create, &bill)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")

		var clientErr *Error
		if errors.As(err, &clientErr) && clientErr.StatusCode == http.StatusBadRequest {
			errorRes := map[string][]string{}
			err := json.Unmarshal([]byte(clientErr.Body), &errorRes)
			if err != nil {
				return nil, res, fmt.Errorf("unmarshaling response body to error response: %w", err)
			}

			if len(errorRes["visit_note"]) > 0 && slices.Contains(errorRes["visit_note"], billExistError) {
				return nil, res, ErrBillExist
			}
		}

		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return bill, res, nil
}

type FindBillsOptions struct {
	*Pagination

	AssignedPhysician []int64   `url:"assigned_physician,omitempty"`
	BillID            []int64   `url:"bill_id,omitempty"`
	FromServiceDate   time.Time `url:"from_service_date,omitempty"`
	ToServiceDate     time.Time `url:"to_service_date,omitempty"`
	Patient           []int64   `url:"patient,omitempty"`
	Practice          []int64   `url:"practice,omitempty"`
	SigningPhysician  []int64   `url:"signing_physician,omitempty"`
	VisitNoteID       []int64   `url:"visit_note_id,omitempty"`
}

func (s *BillService) Find(ctx context.Context, opts *FindBillsOptions) (*Response[[]*Bill], *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "find bills", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &Response[[]*Bill]{}

	res, err := s.client.request(ctx, http.MethodGet, "/bills", opts, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *BillService) Get(ctx context.Context, id int64) (*Bill, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "get bill", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.bill_id", id)))
	defer span.End()

	out := &Bill{}

	res, err := s.client.request(ctx, http.MethodGet, "/bills/"+strconv.FormatInt(id, 10), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

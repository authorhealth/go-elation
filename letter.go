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

type LetterServicer interface {
	Find(ctx context.Context, opts *FindLettersOptions) (*Response[[]*Letter], *http.Response, error)
	Get(ctx context.Context, id int64) (*Letter, *http.Response, error)
}

var _ LetterServicer = (*LetterService)(nil)

type LetterService struct {
	client *Client
}

type Letter struct {
	ID                    int64                `json:"id"`
	Attachments           []*LetterAttachment  `json:"attachments"`            //: [{}],
	Body                  string               `json:"body"`                   //: "Test Body Message",
	DeliveryDate          *time.Time           `json:"delivery_date"`          //: "2021-05-11T21: 16: 23Z",
	DeliveryMethod        string               `json:"delivery_method"`        //: "printed",
	DirectMessageTo       string               `json:"direct_message_to"`      //: "Test",
	DisplayTo             string               `json:"display_to"`             //: "Test Name",
	DocumentDate          time.Time            `json:"document_date"`          //: "2021-05-11T21: 16: 23Z",
	EmailTo               string               `json:"email_to"`               //: "test@email.me",
	FailureUnacknowledged string               `json:"failure_unacknowledged"` //: false,
	FaxStatus             string               `json:"fax_status"`             //: "success",
	FaxAttachments        bool                 `json:"fax_attachments"`        //: true,
	FaxTo                 string               `json:"fax_to"`                 //: "5555555555",
	IsProcessed           bool                 `json:"is_processed"`           //: true,
	LetterType            string               `json:"letter_type"`            //: "provider",
	Patient               int64                `json:"patient"`                //: 140754680086529,
	Practice              int64                `json:"practice"`               //: 140754674450436
	SendToContact         *LetterSendToContact `json:"send_to_contact"`        //: {},
	SendToElationUser     int64                `json:"send_to_elation_user"`   //: 1234,
	SendToName            string               `json:"send_to_name"`           //: "Test Name",
	SignDate              *time.Time           `json:"sign_date"`              //: "2021-05-11T21: 16: 23Z",
	SignedBy              int64                `json:"signed_by"`              //: 12323455,
	Subject               string               `json:"subject"`                //: "Test Subject",
	Tags                  []any                `json:"tags"`                   //: [],
	ToNumber              string               `json:"to_number"`              //: "5555555555",
	ViewedAt              *time.Time           `json:"viewed_at"`              //: "2021-05-11T21: 16: 23Z",
	WithArchive           bool                 `json:"with_archive"`           //: false,
}

type LetterAttachment struct {
	ID           int64  `json:"id"`
	DocumentType string `json:"document_type"`
}

type LetterSendToContact struct {
	ID          int64                           `json:"id"`
	FirstName   string                          `json:"first_name"`
	LastName    string                          `json:"last_name"`
	NPI         string                          `json:"npi"`
	OrgName     string                          `json:"org_name"`
	Specialties []LetterSendToContactSpeciality `json:"specialties"`
}

type LetterSendToContactSpeciality struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type FindLettersOptions struct {
	*Pagination

	DocumentDateGT  time.Time `url:"document_date__gt,omitempty"`
	DocumentDateGTE time.Time `url:"document_date__gte,omitempty"`
	DocumentDateLT  time.Time `url:"document_date__lt,omitempty"`
	DocumentDateLTE time.Time `url:"document_date__lte,omitempty"`
	Patient         int64     `url:"patient,omitempty"`
	Practice        int64     `url:"practice,omitempty"`
	RecipientID     int64     `url:"recipient_id,omitempty"`
	RecipientName   string    `url:"recipient_name,omitempty"`
}

func (s *LetterService) Find(ctx context.Context, opts *FindLettersOptions) (*Response[[]*Letter], *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "find letters", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &Response[[]*Letter]{}

	res, err := s.client.request(ctx, http.MethodGet, "/letters", opts, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *LetterService) Get(ctx context.Context, id int64) (*Letter, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "get letter", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.letter_id", id)))
	defer span.End()

	out := &Letter{}

	res, err := s.client.request(ctx, http.MethodGet, "/letters/"+strconv.FormatInt(id, 10), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

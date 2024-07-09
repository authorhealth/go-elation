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

type MessageThreadServicer interface {
	Find(ctx context.Context, opts *FindMessageThreadsOptions) (*Response[[]*MessageThread], *http.Response, error)
	Get(ctx context.Context, id int64) (*MessageThread, *http.Response, error)
}

var _ MessageThreadServicer = (*MessageThreadService)(nil)

type MessageThreadService struct {
	client *HTTPClient
}

type MessageThread struct {
	CreatedDate  time.Time              `json:"created_date"`  //: "2021-05-05T22:02:10Z",
	DeletedDate  *time.Time             `json:"deleted_date"`  //: null,
	DocumentDate time.Time              `json:"document_date"` //: "2021-05-05T22:02:09Z",
	ChartDate    time.Time              `json:"chart_date"`    //: "2021-05-05T22:02:09Z",
	Patient      int64                  `json:"patient"`       //: 140754786975745,
	Practice     int64                  `json:"practice"`      //: 140754817318916,
	IsUrgent     bool                   `json:"is_urgent"`     //: false,
	ID           int64                  `json:"id"`            //: 140754787631195,
	Members      []MessageThreadMember  `json:"members"`       //: [],
	Messages     []MessageThreadMessage `json:"messages"`      //: []
}

type MessageThreadMember struct {
	ID      int64      `json:"id"`       //: 346292316,
	Thread  int64      `json:"thread"`   //: 346226779,
	User    *int64     `json:"user"`     //: 6,
	Group   *int64     `json:"group"`    //: null,
	Status  string     `json:"status"`   //: "Addressed",
	AckTime *time.Time `json:"ack_time"` //: null
}

type MessageThreadMessage struct {
	Body       string                          `json:"body"`        //: "Patient should have appointment already booked.",
	ID         int64                           `json:"id"`          //: 346423390,
	Patient    int64                           `json:"patient"`     //: 342753281,
	Practice   int64                           `json:"practice"`    //: 655364,
	Thread     int64                           `json:"thread"`      //: 346226779,
	Sender     int64                           `json:"sender"`      //: 6,
	SendDate   time.Time                       `json:"send_date"`   //: "2021-10-19T19:08:48",
	ToDocument *MessageThreadMessageToDocument `json:"to_document"` //: {}
}

type MessageThreadMessageToDocument struct {
	AuthoringPractice int64      `json:"authoring_practice"` //: 65540,
	ChartDate         time.Time  `json:"chart_date"`         //: "2014-01-24T18:13:48Z",
	CreatedDate       time.Time  `json:"created_date"`       //: "2014-01-24T18:13:48Z",
	DeletedDate       *time.Time `json:"deleted_date"`       //: null,
	DocumentDate      time.Time  `json:"document_date"`      //: "2014-01-24T18:59:15Z",
	DocumentType      int64      `json:"document_type"`      //: 91,
	ID                int64      `json:"id"`                 //: 64077889627,
	LastModified      time.Time  `json:"last_modified"`      //: "2016-05-02T18:53:12Z",
	Patient           int64      `json:"patient"`            //: 64072843265,
	SignDate          *time.Time `json:"sign_date"`          //: "2014-01-24T18:13:48Z",
	SignedBy          int64      `json:"signed_by"`          //: 4
}

type FindMessageThreadsOptions struct {
	*Pagination

	Patient         []int64   `url:"patient,omitempty"`
	Practice        []int64   `url:"practice,omitempty"`
	DocumentDateGT  time.Time `url:"document_date__gt,omitempty"`
	DocumentDateGTE time.Time `url:"document_date__gte,omitempty"`
	DocumentDateLT  time.Time `url:"document_date__lt,omitempty"`
	DocumentDateLTE time.Time `url:"document_date__lte,omitempty"`
}

func (s *MessageThreadService) Find(ctx context.Context, opts *FindMessageThreadsOptions) (*Response[[]*MessageThread], *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "find message threads", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &Response[[]*MessageThread]{}

	res, err := s.client.request(ctx, http.MethodGet, "/message_threads", opts, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *MessageThreadService) Get(ctx context.Context, id int64) (*MessageThread, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "get message thread", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.message_thread_id", id)))
	defer span.End()

	out := &MessageThread{}

	res, err := s.client.request(ctx, http.MethodGet, "/message_threads/"+strconv.FormatInt(id, 10), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

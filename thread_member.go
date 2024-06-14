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

type ThreadMemberServicer interface {
	Find(ctx context.Context, opts *FindThreadMembersOptions) (*Response[[]*ThreadMember], *http.Response, error)
	Get(ctx context.Context, id int64) (*ThreadMember, *http.Response, error)
}

var _ ThreadMemberServicer = (*ThreadMemberService)(nil)

type ThreadMemberService struct {
	client *Client
}

type ThreadMember struct {
	ID      int64      `json:"id"`       //: 346292316,
	Thread  int64      `json:"thread"`   //: 346226779,
	User    *int64     `json:"user"`     //: 6,
	Group   *int64     `json:"group"`    //: null,
	Status  string     `json:"status"`   //: "Addressed",
	AckTime *time.Time `json:"ack_time"` //: null
}

type FindThreadMembersOptions struct {
	*Pagination

	Patient  []int64   `url:"patient,omitempty"`
	Practice []int64   `url:"practice,omitempty"`
	User     []int64   `url:"user,omitempty"`
	Group    []int64   `url:"group,omitempty"`
	Status   string    `url:"status,omitempty"`
	AckTime  time.Time `url:"ack_time,omitempty"`
}

func (s *ThreadMemberService) Find(ctx context.Context, opts *FindThreadMembersOptions) (*Response[[]*ThreadMember], *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "find thread members", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &Response[[]*ThreadMember]{}

	res, err := s.client.request(ctx, http.MethodGet, "/thread_members", opts, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *ThreadMemberService) Get(ctx context.Context, id int64) (*ThreadMember, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "get thread member", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.thread_member_id", id)))
	defer span.End()

	out := &ThreadMember{}

	res, err := s.client.request(ctx, http.MethodGet, "/thread_members/"+strconv.FormatInt(id, 10), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

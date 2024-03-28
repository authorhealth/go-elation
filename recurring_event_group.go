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

type RecurringEventGroupServicer interface {
	Create(ctx context.Context, create *RecurringEventGroupCreate) (*RecurringEventGroup, *http.Response, error)
	Find(ctx context.Context, opts *FindRecurringEventGroupsOptions) (*Response[[]*RecurringEventGroup], *http.Response, error)
	Get(ctx context.Context, id int64) (*RecurringEventGroup, *http.Response, error)
	Update(ctx context.Context, id int64, update *RecurringEventGroupUpdate) (*RecurringEventGroup, *http.Response, error)
	Delete(ctx context.Context, id int64) (*http.Response, error)
}

var _ RecurringEventGroupServicer = (*RecurringEventGroupService)(nil)

type RecurringEventGroupService struct {
	client *Client
}

type RecurringEventGroupSchedule struct {
	ID           int64     `json:"id"`
	SeriesStart  string    `json:"series_start"`
	SeriesStop   string    `json:"series_stop"`
	EventTime    string    `json:"event_time"`
	Physician    int64     `json:"physician"`
	Duration     int       `json:"duration"`
	Repeats      string    `json:"repeats"`
	DOWMonday    bool      `json:"dow_monday"`
	DOWTuesday   bool      `json:"dow_tuesday"`
	DOWWednesday bool      `json:"dow_wednesday"`
	DOWThursday  bool      `json:"dow_thursday"`
	DOWFriday    bool      `json:"dow_friday"`
	DOWSaturday  bool      `json:"dow_saturday"`
	DOWSunday    bool      `json:"dow_sunday"`
	Description  string    `json:"description"`
	CreatedDate  time.Time `json:"created_date"`
}

type RecurringEventGroup struct {
	ID       int64 `json:"id"`
	Practice int64 `json:"practice"`

	CreatedDate time.Time  `json:"created_date"`
	DeletedDate *time.Time `json:"deleted_date"`

	Description  string                         `json:"description"`
	Reason       string                         `json:"reason"`
	Schedules    []*RecurringEventGroupSchedule `json:"schedules"`
	TimeSlotType TimeSlotType                   `json:"time_slot_type"`
}

type RecurringEventGroupCreate struct {
	Practice     int64                          `json:"practice"`
	Reason       string                         `json:"reason"`
	Schedules    []*RecurringEventGroupSchedule `json:"schedules"`
	TimeSlotType TimeSlotType                   `json:"time_slot_type"`
}

func (s *RecurringEventGroupService) Create(ctx context.Context, create *RecurringEventGroupCreate) (*RecurringEventGroup, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "create recurring event group", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &RecurringEventGroup{}

	res, err := s.client.request(ctx, http.MethodPost, "/recurring_event_groups", nil, create, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

type FindRecurringEventGroupsOptions struct {
	*Pagination

	Physician    []int64      `url:"physician,omitempty"`
	Practice     []int64      `url:"practice,omitempty"`
	Reason       string       `url:"reason,omitempty"`
	StartDate    time.Time    `url:"start_date,omitempty"`
	EndDate      time.Time    `url:"end_date,omitempty"`
	TimeSlotType TimeSlotType `url:"time_slot_type,omitempty"`
}

func (s *RecurringEventGroupService) Find(ctx context.Context, opts *FindRecurringEventGroupsOptions) (*Response[[]*RecurringEventGroup], *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "find recurring event groups", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &Response[[]*RecurringEventGroup]{}

	res, err := s.client.request(ctx, http.MethodGet, "/recurring_event_groups", opts, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *RecurringEventGroupService) Get(ctx context.Context, id int64) (*RecurringEventGroup, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "get recurring event group", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.recurring_event_group_id", id)))
	defer span.End()

	out := &RecurringEventGroup{}

	res, err := s.client.request(ctx, http.MethodGet, "/recurring_event_groups/"+strconv.FormatInt(id, 10), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

type RecurringEventGroupUpdate struct {
	Reason    string                         `json:"reason"`
	Schedules []*RecurringEventGroupSchedule `json:"schedules"`
}

func (s *RecurringEventGroupService) Update(ctx context.Context, id int64, update *RecurringEventGroupUpdate) (*RecurringEventGroup, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "update recurring event group", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.recurring_event_group_id", id)))
	defer span.End()

	out := &RecurringEventGroup{}

	res, err := s.client.request(ctx, http.MethodPatch, "/recurring_event_groups/"+strconv.FormatInt(id, 10), nil, update, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *RecurringEventGroupService) Delete(ctx context.Context, id int64) (*http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "delete recurring event group", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.recurring_event_group_id", id)))
	defer span.End()

	res, err := s.client.request(ctx, http.MethodDelete, "/recurring_event_groups/"+strconv.FormatInt(id, 10), nil, nil, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return res, fmt.Errorf("making request: %w", err)
	}

	return res, nil
}

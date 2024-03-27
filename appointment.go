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

const (
	AppointmentModeInPerson = "IN_PERSON"
	AppointmentModeVideo    = "VIDEO"
)

type TimeSlotType string

const (
	AppointmentTimeSlotTypeAppointment     TimeSlotType = "appointment"
	AppointmentTimeSlotTypeAppointmentSlot TimeSlotType = "appointment_slot"
	AppointmentTimeSlotTypeEvent           TimeSlotType = "event"
)

type AppointmentServicer interface {
	Create(ctx context.Context, create *AppointmentCreate) (*Appointment, *http.Response, error)
	Find(ctx context.Context, opts *FindAppointmentsOptions) (*Response[[]*Appointment], *http.Response, error)
	Get(ctx context.Context, id int64) (*Appointment, *http.Response, error)
	Update(ctx context.Context, id int64, update *AppointmentUpdate) (*Appointment, *http.Response, error)
	Delete(ctx context.Context, id int64) (*http.Response, error)
}

var _ AppointmentServicer = (*AppointmentService)(nil)

type AppointmentService struct {
	client *Client
}

type AppointmentCreate struct {
	ScheduledDate time.Time `json:"scheduled_date"`
	Reason        string    `json:"reason"`
	Patient       int64     `json:"patient"`
	Physician     int64     `json:"physician"`
	Practice      int64     `json:"practice"`
}

func (s *AppointmentService) Create(ctx context.Context, create *AppointmentCreate) (*Appointment, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "create appointment", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &Appointment{}

	res, err := s.client.request(ctx, http.MethodPost, "/appointments", nil, create, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

type Appointment struct {
	ID                     int64                       `json:"id"`
	ScheduledDate          time.Time                   `json:"scheduled_date"`
	Duration               int                         `json:"duration"`
	TimeSlotType           string                      `json:"time_slot_type"`
	TimeSlotStatus         any                         `json:"time_slot_status"`
	Reason                 string                      `json:"reason"`
	Description            string                      `json:"description"`
	Status                 *AppointmentStatus          `json:"status"`
	ServiceLocation        *AppointmentServiceLocation `json:"service_location"`
	TelehealthDetails      string                      `json:"telehealth_details"`
	Patient                int64                       `json:"patient"`
	Physician              int64                       `json:"physician"`
	Practice               int64                       `json:"practice"`
	RecurringEventSchedule any                         `json:"recurring_event_schedule"`
	BillingDetails         *AppointmentBillingDetails  `json:"billing_details"`
	Payment                *AppointmentPayment         `json:"payment"`
	Metadata               any                         `json:"metadata"`
	CreatedDate            time.Time                   `json:"created_date"`
	LastModifiedDate       time.Time                   `json:"last_modified_date"`
	DeletedDate            *time.Time                  `json:"deleted_date"`
	Mode                   string                      `json:"mode"`
	Instructions           string                      `json:"instructions"`
}

type AppointmentStatus struct {
	Status     string `json:"status"`
	Room       string `json:"room"`
	StatusDate string `json:"status_date"`
}

type AppointmentServiceLocation struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	PlaceOfService int    `json:"place_of_service"`
	AddressLine1   string `json:"address_line1"`
	AddressLine2   string `json:"address_line2"`
	City           string `json:"city"`
	State          string `json:"state"`
	Zip            string `json:"zip"`
	Phone          string `json:"phone"`
}

type AppointmentBillingDetails struct {
	BillingNote            string `json:"billing_note"`
	ReferringProvider      string `json:"referring_provider"`
	ReferringProviderState string `json:"referring_provider_state"`
}

type AppointmentPayment struct {
	ID            int64      `json:"id"`
	Amount        string     `json:"amount"`
	WhenCollected time.Time  `json:"when_collected"`
	Bill          any        `json:"bill"`
	Appointment   int64      `json:"appointment"`
	CreateDate    time.Time  `json:"create_date"`
	DeleteDate    *time.Time `json:"delete_date"`
}

type FindAppointmentsOptions struct {
	*Pagination

	Patient      []int64   `url:"patient,omitempty"`
	Practice     []int64   `url:"practice,omitempty"`
	Physician    []int64   `url:"physician,omitempty"`
	FromDate     time.Time `url:"from_date,omitempty"`
	ToDate       time.Time `url:"to_date,omitempty"`
	TimeSlotType string    `url:"time_slot_type,omitempty"`
}

func (s *AppointmentService) Find(ctx context.Context, opts *FindAppointmentsOptions) (*Response[[]*Appointment], *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "find appointments", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &Response[[]*Appointment]{}

	res, err := s.client.request(ctx, http.MethodGet, "/appointments", opts, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *AppointmentService) Get(ctx context.Context, id int64) (*Appointment, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "get appointment", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.appointment_id", id)))
	defer span.End()

	out := &Appointment{}

	res, err := s.client.request(ctx, http.MethodGet, "/appointments/"+strconv.FormatInt(id, 10), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

type AppointmentUpdate struct {
	Duration          *int                     `json:"duration,omitempty"`
	Instructions      *string                  `json:"instructions,omitempty"`
	Mode              *string                  `json:"mode,omitempty"`
	Status            *AppointmentUpdateStatus `json:"status,omitempty"`
	TelehealthDetails *string                  `json:"telehealth_details,omitempty"`
}

type AppointmentUpdateStatus struct {
	Status string `json:"status"`         // Required
	Room   string `json:"room,omitempty"` // Optional
}

func (s *AppointmentService) Update(ctx context.Context, id int64, update *AppointmentUpdate) (*Appointment, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "update appointment", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.appointment_id", id)))
	defer span.End()

	out := &Appointment{}

	res, err := s.client.request(ctx, http.MethodPatch, "/appointments/"+strconv.FormatInt(id, 10), nil, update, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *AppointmentService) Delete(ctx context.Context, id int64) (*http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "delete appointment", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.appointment_id", id)))
	defer span.End()

	res, err := s.client.request(ctx, http.MethodDelete, "/appointments/"+strconv.FormatInt(id, 10), nil, nil, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return res, fmt.Errorf("making request: %w", err)
	}

	return res, nil
}

package elation

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type AppointmentServicer interface {
	Create(ctx context.Context, appt *AppointmentCreate) (*Appointment, *http.Response, error)
	Find(ctx context.Context, opts *FindAppointmentsOptions) ([]*Appointment, *http.Response, error)
	Get(ctx context.Context, id int64) (*Appointment, *http.Response, error)
	Update(ctx context.Context, id int64, appt *AppointmentUpdate) (*Appointment, *http.Response, error)
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

func (a *AppointmentService) Create(ctx context.Context, appt *AppointmentCreate) (*Appointment, *http.Response, error) {
	out := &Appointment{}

	res, err := a.client.request(ctx, http.MethodPost, "/appointments", nil, appt, &out)
	if err != nil {
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
	Physician              int                         `json:"physician"`
	Practice               int                         `json:"practice"`
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
	Patient      []int64   `url:"patient,omitempty"`
	Practice     []int64   `url:"practice,omitempty"`
	Physician    []int64   `url:"physician,omitempty"`
	FromDate     time.Time `url:"from_date,omitempty"`
	ToDate       time.Time `url:"to_date,omitempty"`
	TimeSlotType string    `url:"time_slot_type,omitempty"`
}

func (a *AppointmentService) Find(ctx context.Context, opts *FindAppointmentsOptions) ([]*Appointment, *http.Response, error) {
	out := &Response[[]*Appointment]{}

	res, err := a.client.request(ctx, http.MethodGet, "/appointments", opts, nil, &out)
	if err != nil {
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out.Results, res, nil
}

func (a *AppointmentService) Get(ctx context.Context, id int64) (*Appointment, *http.Response, error) {
	out := &Appointment{}

	res, err := a.client.request(ctx, http.MethodGet, "/appointments/"+strconv.FormatInt(id, 10), nil, nil, &out)
	if err != nil {
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

type AppointmentUpdate struct {
	Duration          *int    `json:"duration,omitempty"`
	TelehealthDetails *string `json:"telehealth_details,omitempty"`
}

func (a *AppointmentService) Update(ctx context.Context, id int64, appt *AppointmentUpdate) (*Appointment, *http.Response, error) {
	out := &Appointment{}

	res, err := a.client.request(ctx, http.MethodPatch, "/appointments/"+strconv.FormatInt(id, 10), nil, appt, &out)
	if err != nil {
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (a *AppointmentService) Delete(ctx context.Context, id int64) (*http.Response, error) {
	res, err := a.client.request(ctx, http.MethodDelete, "/appointments/"+strconv.FormatInt(id, 10), nil, nil, nil)
	if err != nil {
		return res, fmt.Errorf("making request: %w", err)
	}

	return res, nil
}

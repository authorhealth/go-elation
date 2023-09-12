package elation

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAppointmentService_Create(t *testing.T) {
	assert := assert.New(t)

	expected := &AppointmentCreate{
		ScheduledDate: time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
		Reason:        "reason",
		Patient:       1,
		Physician:     2,
		Practice:      3,
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodPost, r.Method)
		assert.Equal("/appointments", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		assert.NoError(err)

		actual := &AppointmentCreate{}
		err = json.Unmarshal(body, actual)
		assert.NoError(err)

		assert.Equal(expected.ScheduledDate, actual.ScheduledDate)
		assert.Equal(expected.Reason, actual.Reason)
		assert.Equal(expected.Patient, actual.Patient)
		assert.Equal(expected.Physician, actual.Physician)
		assert.Equal(expected.Practice, actual.Practice)

		b, err := json.Marshal(&Appointment{})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := AppointmentService{client}

	created, res, err := svc.Create(context.Background(), expected)
	assert.NotNil(created)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestAppointmentService_Find(t *testing.T) {
	assert := assert.New(t)

	opts := &FindAppointmentsOptions{
		Pagination: &Pagination{
			Limit:  1,
			Offset: 2,
		},

		Patient:      []int64{1},
		Practice:     []int64{2},
		Physician:    []int64{3},
		FromDate:     time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
		ToDate:       time.Date(2023, 5, 20, 0, 0, 0, 0, time.UTC),
		TimeSlotType: "time slot type",
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/appointments", r.URL.Path)

		patient := r.URL.Query()["patient"]
		practice := r.URL.Query()["practice"]
		physician := r.URL.Query()["physician"]
		fromDate := r.URL.Query().Get("from_date")
		toDate := r.URL.Query().Get("to_date")
		timeSlotType := r.URL.Query().Get("time_slot_type")

		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")

		assert.Equal(opts.Patient, sliceStrToInt64(patient))
		assert.Equal(opts.Practice, sliceStrToInt64(practice))
		assert.Equal(opts.Physician, sliceStrToInt64(physician))
		assert.Equal(opts.FromDate.Format(time.RFC3339), fromDate)
		assert.Equal(opts.ToDate.Format(time.RFC3339), toDate)
		assert.Equal(opts.TimeSlotType, timeSlotType)

		assert.Equal(opts.Pagination.Limit, strToInt(limit))
		assert.Equal(opts.Pagination.Offset, strToInt(offset))

		b, err := json.Marshal(Response[[]*Appointment]{
			Results: []*Appointment{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			},
		})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := AppointmentService{client}

	found, res, err := svc.Find(context.Background(), opts)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestAppointmentService_Get(t *testing.T) {
	assert := assert.New(t)

	var id int64 = 1

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/appointments/"+strconv.FormatInt(id, 10), r.URL.Path)

		b, err := json.Marshal(&Appointment{
			ID: id,
		})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := AppointmentService{client}

	found, res, err := svc.Get(context.Background(), id)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestAppointmentService_Update(t *testing.T) {
	assert := assert.New(t)

	var id int64 = 1
	expected := &AppointmentUpdate{
		Duration:          Ptr(30),
		TelehealthDetails: Ptr("telehealth details"),
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodPatch, r.Method)
		assert.Equal("/appointments/"+strconv.FormatInt(id, 10), r.URL.Path)

		body, err := io.ReadAll(r.Body)
		assert.NoError(err)

		actual := &AppointmentUpdate{}
		err = json.Unmarshal(body, actual)
		assert.NoError(err)

		assert.Equal(expected, actual)

		b, err := json.Marshal(&Appointment{})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := AppointmentService{client}

	found, res, err := svc.Update(context.Background(), id, expected)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestAppointmentService_Delete(t *testing.T) {
	assert := assert.New(t)

	var id int64 = 1

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodDelete, r.Method)
		assert.Equal("/appointments/"+strconv.FormatInt(id, 10), r.URL.Path)
	}))
	defer srv.Close()

	client := NewClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := AppointmentService{client}

	res, err := svc.Delete(context.Background(), id)
	assert.NotNil(res)
	assert.NoError(err)
}

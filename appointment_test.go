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

	expectedAppt := &AppointmentCreate{
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

		body, err := io.ReadAll(r.Body)
		assert.NoError(err)

		actualAppt := &AppointmentCreate{}
		err = json.Unmarshal(body, actualAppt)
		assert.NoError(err)

		assert.Equal(expectedAppt.ScheduledDate, actualAppt.ScheduledDate)
		assert.Equal(expectedAppt.Reason, actualAppt.Reason)
		assert.Equal(expectedAppt.Patient, actualAppt.Patient)
		assert.Equal(expectedAppt.Physician, actualAppt.Physician)
		assert.Equal(expectedAppt.Practice, actualAppt.Practice)

		b, err := json.Marshal(&Appointment{})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}))
	defer srv.Close()

	client := NewClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := AppointmentService{client}

	appt, res, err := svc.Create(context.Background(), expectedAppt)
	assert.NotNil(appt)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestAppointmentService_Find(t *testing.T) {
	assert := assert.New(t)

	opts := &FindAppointmentsOptions{
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

		patient := r.URL.Query()["patient"]
		practice := r.URL.Query()["practice"]
		physician := r.URL.Query()["physician"]
		fromDate := r.URL.Query().Get("from_date")
		toDate := r.URL.Query().Get("to_date")
		timeSlotType := r.URL.Query().Get("time_slot_type")

		assert.Equal(opts.Patient, sliceStrToInt64(patient))
		assert.Equal(opts.Practice, sliceStrToInt64(practice))
		assert.Equal(opts.Physician, sliceStrToInt64(physician))
		assert.Equal(opts.FromDate.Format(time.RFC3339), fromDate)
		assert.Equal(opts.ToDate.Format(time.RFC3339), toDate)
		assert.Equal(opts.TimeSlotType, timeSlotType)

		b, err := json.Marshal(Response[[]*Appointment]{
			Count:    0,
			Next:     "",
			Previous: "",
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
		w.Write(b)
	}))
	defer srv.Close()

	client := NewClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := AppointmentService{client}

	appts, res, err := svc.Find(context.Background(), opts)
	assert.NotNil(appts)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestAppointmentService_Get(t *testing.T) {
	assert := assert.New(t)

	var apptID int64 = 1

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/appointments/"+strconv.FormatInt(apptID, 10), r.URL.Path)

		b, err := json.Marshal(&Appointment{
			ID: apptID,
		})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}))
	defer srv.Close()

	client := NewClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := AppointmentService{client}

	appt, res, err := svc.Get(context.Background(), apptID)
	assert.NotNil(appt)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestAppointmentService_Update(t *testing.T) {
	assert := assert.New(t)

	var apptID int64 = 1
	expectedAppt := &AppointmentUpdate{
		Duration:          Ptr(30),
		TelehealthDetails: Ptr("telehealth details"),
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodPatch, r.Method)
		assert.Equal("/appointments/"+strconv.FormatInt(apptID, 10), r.URL.Path)

		body, err := io.ReadAll(r.Body)
		assert.NoError(err)

		actualAppt := &AppointmentUpdate{}
		err = json.Unmarshal(body, actualAppt)
		assert.NoError(err)

		assert.Equal(expectedAppt.Duration, actualAppt.Duration)
		assert.Equal(expectedAppt.TelehealthDetails, actualAppt.TelehealthDetails)

		b, err := json.Marshal(&Appointment{})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}))
	defer srv.Close()

	client := NewClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := AppointmentService{client}

	appt, res, err := svc.Update(context.Background(), apptID, expectedAppt)
	assert.NotNil(appt)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestAppointmentService_Delete(t *testing.T) {
	assert := assert.New(t)

	var apptID int64 = 1

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodDelete, r.Method)
		assert.Equal("/appointments/"+strconv.FormatInt(apptID, 10), r.URL.Path)
	}))
	defer srv.Close()

	client := NewClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := AppointmentService{client}

	res, err := svc.Delete(context.Background(), apptID)
	assert.NotNil(res)
	assert.NoError(err)
}

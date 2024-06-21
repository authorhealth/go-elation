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

func TestRecurringEventGroupService_Create(t *testing.T) {
	assert := assert.New(t)

	expected := &RecurringEventGroupCreate{
		Practice: 1,
		Reason:   "Appt follow-up",
		Schedules: []*RecurringEventGroupSchedule{
			{
				ID:           5,
				SeriesStart:  "2024-01-01",
				SeriesStop:   "2024-02-01",
				EventTime:    "10:30:00",
				Physician:    2,
				Duration:     15,
				Repeats:      "Weekly",
				DOWMonday:    true,
				DOWTuesday:   false,
				DOWWednesday: false,
				DOWThursday:  false,
				DOWFriday:    false,
				DOWSaturday:  false,
				DOWSunday:    false,
				Description:  "Appt follow-up",
				CreatedDate:  time.Date(2024, 3, 27, 10, 30, 0, 0, time.UTC),
			},
		},
		TimeSlotType: AppointmentTimeSlotTypeAppointmentSlot,
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodPost, r.Method)
		assert.Equal("/recurring_event_groups", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		assert.NoError(err)

		actual := &RecurringEventGroupCreate{}
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

	client := NewHttpClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := RecurringEventGroupService{client}

	created, res, err := svc.Create(context.Background(), expected)
	assert.NotNil(created)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestRecurringEventGroupService_Find(t *testing.T) {
	assert := assert.New(t)

	opts := &FindRecurringEventGroupsOptions{
		Pagination: &Pagination{
			Limit:  1,
			Offset: 2,
		},

		Physician:    []int64{1},
		Practice:     []int64{2},
		Reason:       "Some reason",
		TimeSlotType: AppointmentTimeSlotTypeEvent,
		StartDate:    "2024-01-01",
		EndDate:      "2024-03-01",
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/recurring_event_groups", r.URL.Path)

		practice := r.URL.Query()["practice"]
		physician := r.URL.Query()["physician"]
		reason := r.URL.Query().Get("reason")
		timeSlotType := r.URL.Query().Get("time_slot_type")
		startDate := r.URL.Query().Get("start_date")
		endDate := r.URL.Query().Get("end_date")

		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")

		assert.Equal(opts.Practice, commaStrToInt64(practice[0]))
		assert.Equal(opts.Physician, commaStrToInt64(physician[0]))
		assert.Equal(opts.Reason, reason)
		assert.Equal(opts.TimeSlotType, TimeSlotType(timeSlotType))
		assert.Equal(opts.StartDate, startDate)
		assert.Equal(opts.EndDate, endDate)

		assert.Equal(opts.Pagination.Limit, strToInt(limit))
		assert.Equal(opts.Pagination.Offset, strToInt(offset))

		b, err := json.Marshal(Response[[]*RecurringEventGroup]{
			Results: []*RecurringEventGroup{
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

	client := NewHttpClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := RecurringEventGroupService{client}

	found, res, err := svc.Find(context.Background(), opts)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestRecurringEventGroupService_Find_Multiple_Params(t *testing.T) {
	assert := assert.New(t)

	opts := &FindRecurringEventGroupsOptions{
		Pagination: &Pagination{
			Limit:  1,
			Offset: 2,
		},

		Physician:    []int64{1, 2, 3},
		Practice:     []int64{7, 8, 9},
		Reason:       "Some reason",
		TimeSlotType: AppointmentTimeSlotTypeEvent,
		StartDate:    "2024-01-01",
		EndDate:      "2024-03-01",
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/recurring_event_groups", r.URL.Path)

		practice := r.URL.Query()["practice"]
		physician := r.URL.Query()["physician"]
		reason := r.URL.Query().Get("reason")
		timeSlotType := r.URL.Query().Get("time_slot_type")
		startDate := r.URL.Query().Get("start_date")
		endDate := r.URL.Query().Get("end_date")

		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")

		assert.Equal(opts.Practice, commaStrToInt64(practice[0]))
		assert.Equal(opts.Physician, commaStrToInt64(physician[0]))
		assert.Equal(opts.Reason, reason)
		assert.Equal(opts.TimeSlotType, TimeSlotType(timeSlotType))
		assert.Equal(opts.StartDate, startDate)
		assert.Equal(opts.EndDate, endDate)

		assert.Equal(opts.Pagination.Limit, strToInt(limit))
		assert.Equal(opts.Pagination.Offset, strToInt(offset))

		b, err := json.Marshal(Response[[]*RecurringEventGroup]{
			Results: []*RecurringEventGroup{
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

	client := NewHttpClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := RecurringEventGroupService{client}

	found, res, err := svc.Find(context.Background(), opts)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestRecurringEventGroupService_Get(t *testing.T) {
	assert := assert.New(t)

	var id int64 = 1

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/recurring_event_groups/"+strconv.FormatInt(id, 10), r.URL.Path)

		b, err := json.Marshal(&RecurringEventGroup{
			ID: id,
		})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewHttpClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := RecurringEventGroupService{client}

	found, res, err := svc.Get(context.Background(), id)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestRecurringEventGroupService_Update(t *testing.T) {
	assert := assert.New(t)

	var id int64 = 1
	expected := &RecurringEventGroupUpdate{
		Reason: "Appt follow-up",
		Schedules: []*RecurringEventGroupSchedule{
			{
				ID:           5,
				SeriesStart:  "2024-01-01",
				SeriesStop:   "2024-03-01",
				EventTime:    "11:00:00",
				Physician:    2,
				Duration:     15,
				Repeats:      "Weekly",
				DOWMonday:    true,
				DOWTuesday:   true,
				DOWWednesday: false,
				DOWThursday:  false,
				DOWFriday:    false,
				DOWSaturday:  false,
				DOWSunday:    false,
				Description:  "Appt follow-up",
				CreatedDate:  time.Date(2024, 3, 27, 10, 30, 0, 0, time.UTC),
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodPatch, r.Method)
		assert.Equal("/recurring_event_groups/"+strconv.FormatInt(id, 10), r.URL.Path)

		body, err := io.ReadAll(r.Body)
		assert.NoError(err)

		actual := &RecurringEventGroupUpdate{}
		err = json.Unmarshal(body, actual)
		assert.NoError(err)

		assert.Equal(expected, actual)

		b, err := json.Marshal(&RecurringEventGroup{})
		assert.NoError(err)

		w.Header().Set("Content-Type", "application/json")
		//nolint
		w.Write(b)
	}))
	defer srv.Close()

	client := NewHttpClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := RecurringEventGroupService{client}

	found, res, err := svc.Update(context.Background(), id, expected)
	assert.NotNil(found)
	assert.NotNil(res)
	assert.NoError(err)
}

func TestRecurringEventGroupService_Delete(t *testing.T) {
	assert := assert.New(t)

	var id int64 = 1

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenRequest(w, r) {
			return
		}

		assert.Equal(http.MethodDelete, r.Method)
		assert.Equal("/recurring_event_groups/"+strconv.FormatInt(id, 10), r.URL.Path)
	}))
	defer srv.Close()

	client := NewHttpClient(srv.Client(), srv.URL+"/token", "", "", srv.URL)
	svc := RecurringEventGroupService{client}

	res, err := svc.Delete(context.Background(), id)
	assert.NotNil(res)
	assert.NoError(err)
}

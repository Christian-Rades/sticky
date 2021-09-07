package calendar

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestEventInsertion(t *testing.T) {
	now := time.Now()
	past := now.Add(-5 * time.Minute)
	future := now.Add(time.Minute)
	busyDate := now.Add(24 * time.Hour)

	calendar := &Calendar{
		Events: Timeline{[]*Event{
			{
				Id:       uuid.NewString(),
				Title:    "busy1",
				Start:    busyDate,
				Duration: time.Hour,
			},
			{
				Id:       uuid.NewString(),
				Title:    "busy2",
				Start:    busyDate.Add(time.Hour),
				Duration: time.Hour,
			},
		}},
	}

	cases := map[string]struct {
		ev  Event
		err error
	}{
		"Valid Event with Id": {
			Event{
				Id:       uuid.NewString(),
				Title:    "title1",
				Start:    future,
				Duration: time.Hour,
			},
			nil,
		},
		"Valid Event without Id": {
			Event{
				Id:       "",
				Title:    "title2",
				Start:    future,
				Duration: time.Hour,
			},
			nil,
		},
		"Missing title": {
			Event{
				Id:       "",
				Title:    "",
				Start:    future,
				Duration: time.Hour,
			},
			MissingTitle,
		},
		"Start date in the past": {
			Event{
				Id:       "",
				Title:    "title3",
				Start:    past,
				Duration: time.Hour,
			},
			EventInThePast{past},
		},
		"Duration zero": {
			Event{
				Id:       "",
				Title:    "title4",
				Start:    future,
				Duration: 0,
			},
			InvalidDuration{0},
		},
		"Duration negative": {
			Event{
				Id:       "",
				Title:    "title5",
				Start:    future,
				Duration: -5 * time.Minute,
			},
			InvalidDuration{-5 * time.Minute},
		},
	}
	for test, data := range cases {
		t.Log(test)
		err := calendar.NewEvent(&data.ev)
		if data.err == nil && err != nil {
			t.Errorf("unexpected error\n\tgot: %q", err)
			continue
		}

		if data.err != nil && !errors.Is(err, data.err) {
			t.Errorf("unexpected error\n\tgot: %q\n\twanted: %q", err, data.err)
		}
	}
}

func TestSchedulingConflicts(t *testing.T) {
	now := time.Now()
	busyDate := now.Add(24 * time.Hour)
	previousEvent := &Event{
		Id:       uuid.NewString(),
		Title:    "busy1",
		Start:    busyDate,
		Duration: time.Hour,
	}

	calendar := &Calendar{
		Events: Timeline{[]*Event{
			previousEvent,
		}},
	}
	cases := map[string]struct {
		ev  Event
		err error
	}{
		"Collision start during event": {
			Event{
				Id:       "",
				Title:    "title6",
				Start:    previousEvent.Start.Add(10 * time.Minute),
				Duration: time.Hour,
			},
			SchedulingConflict{
				NewEvent: &Event{
					Title: "title6",
					Start: busyDate.Add(10 * time.Minute),
				},
				ExistingEvent: previousEvent,
			},
		},
		"Collision running into event": {
			Event{
				Id:       "",
				Title:    "title7",
				Start:    previousEvent.Start.Add(-10 * time.Minute),
				Duration: time.Hour,
			},
			SchedulingConflict{
				NewEvent: &Event{
					Title: "title7",
					Start: busyDate.Add(-10 * time.Minute),
				},
				ExistingEvent: previousEvent,
			},
		},
	}
	for test, data := range cases {
		t.Log(test)
		err := calendar.NewEvent(&data.ev)
		if err.Error() != data.err.Error() {
			t.Errorf("unexpected error\n\tgot: %q\n\twanted: %q", err, data.err)
		}
	}
}

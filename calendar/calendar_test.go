package calendar

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewEvent(t *testing.T) {
	c := &Calendar{}

	title := "my event"
	start := time.Date(2021, 8, 26, 0, 0, 0, 0, time.Local)
	duration, _ := time.ParseDuration("45m")
	ev := &Event{
		Id:       "",
		Title:    title,
		Start:    start,
		Duration: duration,
	}

	err := c.NewEvent(ev)

	if err != nil {
		t.Fatalf("unexpected error %q", err)
	}
}

func TestEventInsertion(t *testing.T) {
	calendar := &Calendar{}
	cases := map[string]struct {
		ev  Event
		err error
	}{
		"Valid Event with Id": {
			Event{
				Id:       uuid.NewString(),
				Title:    "title1",
				Start:    time.Now().Add(time.Minute),
				Duration: time.Hour,
			},
			nil,
		},
		"Valid Event without Id": {
			Event{
				Id:       "",
				Title:    "title2",
				Start:    time.Now().Add(time.Minute),
				Duration: time.Hour,
			},
			nil,
		},
		"Missing title": {
			Event{
				Id:       "",
				Title:    "",
				Start:    time.Now().Add(time.Minute),
				Duration: time.Hour,
			},
			MissingTitle,
		},
		"Start date in the past": {
			Event{
				Id:       "",
				Title:    "title3",
				Start:    time.Now().Add(-5 * time.Minute),
				Duration: time.Hour,
			},
			EventInThePast{},
		},
		"Duration zero": {
			Event{
				Id:       "",
				Title:    "title4",
				Start:    time.Now().Add(time.Minute),
				Duration: 0,
			},
			InvalidDruation{},
		},
		"Duration negative": {
			Event{
				Id:       "",
				Title:    "title5",
				Start:    time.Now().Add(time.Minute),
				Duration: -5 * time.Minute,
			},
			InvalidDruation{},
		},
	}
	for test, data := range cases {
		t.Log(test)
		err := calendar.NewEvent(&data.ev)
		if data.err == nil && err != nil {
			t.Errorf("unexpected error\n\tgot: %q", err)
			continue
		}
		if data.err != nil && !errors.As(err, &data.err) {
			t.Errorf("unexpected error\n\tgot: %T\n\twanted: %T", err, data.err)
		}
	}
}

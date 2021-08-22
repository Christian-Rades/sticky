package calendar

import (
	"errors"
	"fmt"
	"time"
)

type EventInThePast struct {
	Start time.Time
}

func (e EventInThePast) Error() string {
	return fmt.Sprintf(
		"an event start date cannot be in the past: %q",
		e.Start.Format(time.RFC1123),
	)
}

var MissingTitle = errors.New("event is missing a title")

type InvalidDruation struct {
	Duration time.Duration
}

func (e InvalidDruation) Error() string {
	return fmt.Sprintf(
		"a duration must positive and non zero: %q",
		e.Duration.String(),
	)
}

type Calendar struct {
	Id     string
	Events []Event
	// will probably become necessary
	// from time.Time
	// to time.Time
}

type Event struct {
	Id       string
	Title    string
	Start    time.Time
	Duration time.Duration
}

func (c *Calendar) NewEvent(ev *Event) error {
	if len(ev.Title) == 0 {
		return MissingTitle
	}
	if !ev.Start.After(time.Now()) {
		return EventInThePast{ev.Start}
	}
	fmt.Println("new event created")
	return nil
}

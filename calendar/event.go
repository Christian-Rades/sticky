package calendar

import (
	"errors"
	"fmt"
	"time"
)

type Event struct {
	Id       string
	Title    string
	Start    time.Time
	Duration time.Duration
}

type EventInThePast struct {
	Start time.Time
}

type InvalidDuration struct {
	Duration time.Duration
}

var MissingTitle = errors.New("event is missing a title")

func (e EventInThePast) Error() string {
	return fmt.Sprintf(
		"an event start date cannot be in the past: %q",
		e.Start.Format(time.RFC1123),
	)
}

func (e InvalidDuration) Error() string {
	return fmt.Sprintf(
		"a duration must positive and non zero: %q",
		e.Duration.String(),
	)
}

func (e *Event) validate() error {
	if len(e.Title) == 0 {
		return MissingTitle
	}
	if !e.Start.After(time.Now()) {
		return EventInThePast{e.Start}
	}
	if e.Duration.Seconds() <= 0 {
		return InvalidDuration{e.Duration}
	}
	return nil
}



func (e *Event) collides(other *Event) bool {
	return e.contains(other.Start) || e.contains(other.end())
}

func (e *Event) contains(t time.Time) bool {
	return e.Start.Before(t) && t.Before(e.end())
}

func (e *Event) end() time.Time {
	return e.Start.Add(e.Duration)
}

type Timeline struct { Events []*Event }

func (t Timeline) Len() int {
	return len(t.Events)
}

func (t Timeline) Swap(i,j int) {
	t.Events[i], t.Events[j] = t.Events[j], t.Events[i]
}

func (t Timeline) Less(i,j int) bool {
	return t.Events[j].Start.After(t.Events[i].Start)
}

func (t Timeline) findOverlap(ev *Event) []*Event {
	var collisions []*Event
	// TODO optimize:
	// Set Upper bound to first event starting after ev.end()
	for _,e := range t.Events {
		if e.collides(ev) {
			collisions = append(collisions, e)
		}
	}
	return collisions
}

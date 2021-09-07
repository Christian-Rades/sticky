package calendar

import (
	"fmt"
	"time"
)

type Calendar struct {
	Id     string
	Events Timeline
	// will probably become necessary
	// from time.Time
	// to time.Time
}

type SchedulingConflict struct {
	NewEvent *Event
	ExistingEvent *Event
}

func (sc SchedulingConflict) Error() string {
	return fmt.Sprintf(
		"scheduling conflict between %10q on %s \n and \n %10q on %s ",
		sc.NewEvent.Title,
		sc.NewEvent.Start.Format(time.RFC1123),
		sc.ExistingEvent.Title,
		sc.ExistingEvent.Start.Format(time.RFC1123),
	)
}

func (c *Calendar) NewEvent(ev *Event) error {
	fmt.Println("new event created")
	if err := ev.validate(); err != nil {
		return err
	}
	if collisions := c.Events.findOverlap(ev); len(collisions) > 0 {
		return SchedulingConflict{
			NewEvent:       ev,
			ExistingEvent: collisions[0],
		}
	}
	return nil
}

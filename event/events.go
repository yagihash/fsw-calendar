package event

import "google.golang.org/api/calendar/v3"

type Events []*Event

func NewEvents(items []*calendar.Event) Events {
	events := make([]*Event, len(items))

	for i := 0; i < len(items); i++ {
		events[i] = &Event{items[i]}
	}

	return events
}

func (es Events) Diff(another Events) (negative, positive Events) {
	for _, e := range another {
		if !(es.Has(e) || negative.Has(e)) {
			negative = append(negative, e)
		}
	}

	for _, e := range es {
		if !(another.Has(e) || positive.Has(e)) {
			positive = append(positive, e)
		}
	}

	return
}

func (es Events) Has(b *Event) bool {
	for _, a := range es {
		if a.Equals(b) {
			return true
		}
	}

	return false
}

func (es Events) Unique() (unique Events) {
	for i, e := range es {
		// note: used in the case that the original calendar is broken. no need to ensure uniqueness seriously.
		if es[i+1:].Has(e) {
			// do nothing
		} else {
			unique = append(unique, e)
		}
	}

	return
}

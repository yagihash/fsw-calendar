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

func (es Events) Diff(another Events) (negative, positive Events, err error) {
	for _, e := range another {
		has, err := es.Has(e)
		if err != nil {
			return nil, nil, err
		}
		negHas, err := negative.Has(e)
		if err != nil {
			return nil, nil, err
		}
		if !(has || negHas) {
			negative = append(negative, e)
		}
	}

	for _, e := range es {
		has, err := another.Has(e)
		if err != nil {
			return nil, nil, err
		}
		posHas, err := positive.Has(e)
		if err != nil {
			return nil, nil, err
		}
		if !(has || posHas) {
			positive = append(positive, e)
		}
	}

	return
}

func (es Events) Has(b *Event) (bool, error) {
	for _, a := range es {
		ok, err := a.Equals(b)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}

	return false, nil
}

func (es Events) Unique() (Events, error) {
	var unique Events
	for i, e := range es {
		// note: used in the case that the original calendar is broken. no need to ensure uniqueness seriously.
		has, err := es[i+1:].Has(e)
		if err != nil {
			return nil, err
		}
		if !has {
			unique = append(unique, e)
		}
	}

	return unique, nil
}

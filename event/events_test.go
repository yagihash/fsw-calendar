package event

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"google.golang.org/api/calendar/v3"
)

var (
	A = &calendar.Event{
		Summary: "A",
		Start:   &calendar.EventDateTime{DateTime: time.Now().Format(time.RFC3339)},
		End:     &calendar.EventDateTime{DateTime: time.Now().Add(25 * time.Minute).Format(time.RFC3339)},
	}
	B = &calendar.Event{
		Summary: "B",
		Start:   &calendar.EventDateTime{DateTime: time.Now().Add(60 * time.Minute).Format(time.RFC3339)},
		End:     &calendar.EventDateTime{DateTime: time.Now().Add(85 * time.Minute).Format(time.RFC3339)},
	}
	C = &calendar.Event{
		Summary: "C",
		Start:   &calendar.EventDateTime{DateTime: time.Now().Add(120 * time.Minute).Format(time.RFC3339)},
		End:     &calendar.EventDateTime{DateTime: time.Now().Add(145 * time.Minute).Format(time.RFC3339)},
	}
	D = &calendar.Event{
		Summary: "D",
		Start:   &calendar.EventDateTime{DateTime: time.Now().Add(180 * time.Minute).Format(time.RFC3339)},
		End:     &calendar.EventDateTime{DateTime: time.Now().Add(205 * time.Minute).Format(time.RFC3339)},
	}
)

func TestEvents_Diff(t *testing.T) {
	one := Events{A, B, C}
	another := Events{A, B, D}

	wantNegative := Events{D}
	wantPositive := Events{C}

	gotNegative, gotPositive := one.Diff(another)

	if diff := cmp.Diff(gotNegative, wantNegative); diff != "" {
		t.Errorf("got an unexpected diff:\n%s", diff)
	}

	if diff := cmp.Diff(gotPositive, wantPositive); diff != "" {
		t.Errorf("got an unexpected diff:\n%s", diff)
	}
}

func TestEvents_Has(t *testing.T) {
	events := Events{A}

	if !events.Has(A) {
		t.Errorf("Events.Has(A) should return true but got false")
	}

	if events.Has(B) {
		t.Errorf("Events.Has(B) should return false but got true")
	}
}

func TestFetch(t *testing.T) {
	// TODO: change implementation of Fetch so that you can test easily
}

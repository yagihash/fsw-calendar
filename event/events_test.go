package event

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"google.golang.org/api/calendar/v3"
)

var (
	A = &Event{
		&calendar.Event{
			Summary: "A",
			Start:   &calendar.EventDateTime{DateTime: time.Now().Format(time.RFC3339)},
			End:     &calendar.EventDateTime{DateTime: time.Now().Add(25 * time.Minute).Format(time.RFC3339)},
		},
	}
	B = &Event{
		&calendar.Event{
			Summary: "B",
			Start:   &calendar.EventDateTime{DateTime: time.Now().Add(60 * time.Minute).Format(time.RFC3339)},
			End:     &calendar.EventDateTime{DateTime: time.Now().Add(85 * time.Minute).Format(time.RFC3339)},
		},
	}
	C = &Event{
		&calendar.Event{
			Summary: "C",
			Start:   &calendar.EventDateTime{DateTime: time.Now().Add(120 * time.Minute).Format(time.RFC3339)},
			End:     &calendar.EventDateTime{DateTime: time.Now().Add(145 * time.Minute).Format(time.RFC3339)},
		},
	}
	D = &Event{
		&calendar.Event{
			Summary: "D",
			Start:   &calendar.EventDateTime{DateTime: time.Now().Add(180 * time.Minute).Format(time.RFC3339)},
			End:     &calendar.EventDateTime{DateTime: time.Now().Add(205 * time.Minute).Format(time.RFC3339)},
		},
	}

	eventBrokenStartDateTime = &Event{
		&calendar.Event{
			Summary: "eventBrokenDateTime",
			Start:   &calendar.EventDateTime{DateTime: "---"},
			End:     &calendar.EventDateTime{DateTime: time.Now().Add(25 * time.Minute).Format(time.RFC3339)},
		},
	}
	eventBrokenEndDateTime = &Event{
		&calendar.Event{
			Summary: "eventBrokenDateTime",
			Start:   &calendar.EventDateTime{DateTime: time.Now().Format(time.RFC3339)},
			End:     &calendar.EventDateTime{DateTime: "---"},
		},
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
	type args struct {
		b *Event
	}

	tests := []struct {
		name        string
		events      Events
		args        args
		want        bool
		expectPanic bool
	}{
		{
			name:   "True",
			events: Events{A},
			args: args{
				b: A,
			},
			want:        true,
			expectPanic: false,
		},
		{
			name:   "False",
			events: Events{A},
			args: args{
				b: B,
			},
			want:        false,
			expectPanic: false,
		},
		{
			name:   "BrokenStartDatetimeInEvents",
			events: Events{eventBrokenStartDateTime},
			args: args{
				b: A,
			},
			want:        false,
			expectPanic: true,
		},
		{
			name:   "BrokenEndDatetimeInEvents",
			events: Events{eventBrokenEndDateTime},
			args: args{
				b: A,
			},
			want:        false,
			expectPanic: true,
		},
		{
			name:   "BrokenStartDatetimeInB",
			events: Events{A},
			args: args{
				b: eventBrokenStartDateTime,
			},
			want:        false,
			expectPanic: true,
		},
		{
			name:   "BrokenEndDatetimeInB",
			events: Events{A},
			args: args{
				b: eventBrokenEndDateTime,
			},
			want:        false,
			expectPanic: true,
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			defer func() {
				err := recover()

				switch {
				case err != nil && c.expectPanic:
					// OK
				case err != nil && !c.expectPanic:
					t.Errorf("unexpected panic: %v", err)
				case err == nil && c.expectPanic:
					t.Error("expected panic but did not panic")
				case err == nil && !c.expectPanic:
					// OK
				}
			}()

			if diff := cmp.Diff(c.events.Has(c.args.b), c.want); diff != "" {
				t.Errorf("got an unexpected diff:\n%s", diff)
			}
		})
	}
}

func TestEvents_Unique(t *testing.T) {
	redundant := Events{A, A, B, B, C, C}
	want := Events{A, B, C}

	got := redundant.Unique()

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("got an unexpected diff:\n%s", diff)
	}
}

func TestFetch(t *testing.T) {
	// TODO: change implementation of Fetch so that you can test easily
}

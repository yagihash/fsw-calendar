package event

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"google.golang.org/api/calendar/v3"
)

var (
	jst = time.FixedZone("JST", 9*60*60)
	t0  = time.Date(2025, 1, 1, 9, 0, 0, 0, jst)

	A = &Event{
		&calendar.Event{
			Summary: "A",
			Start:   &calendar.EventDateTime{DateTime: t0.Format(time.RFC3339)},
			End:     &calendar.EventDateTime{DateTime: t0.Add(25 * time.Minute).Format(time.RFC3339)},
		},
	}
	B = &Event{
		&calendar.Event{
			Summary: "B",
			Start:   &calendar.EventDateTime{DateTime: t0.Add(60 * time.Minute).Format(time.RFC3339)},
			End:     &calendar.EventDateTime{DateTime: t0.Add(85 * time.Minute).Format(time.RFC3339)},
		},
	}
	C = &Event{
		&calendar.Event{
			Summary: "C",
			Start:   &calendar.EventDateTime{DateTime: t0.Add(120 * time.Minute).Format(time.RFC3339)},
			End:     &calendar.EventDateTime{DateTime: t0.Add(145 * time.Minute).Format(time.RFC3339)},
		},
	}
	D = &Event{
		&calendar.Event{
			Summary: "D",
			Start:   &calendar.EventDateTime{DateTime: t0.Add(180 * time.Minute).Format(time.RFC3339)},
			End:     &calendar.EventDateTime{DateTime: t0.Add(205 * time.Minute).Format(time.RFC3339)},
		},
	}

	eventBrokenStartDateTime = &Event{
		&calendar.Event{
			Summary: "eventBrokenDateTime",
			Start:   &calendar.EventDateTime{DateTime: "---"},
			End:     &calendar.EventDateTime{DateTime: t0.Add(25 * time.Minute).Format(time.RFC3339)},
		},
	}
	eventBrokenEndDateTime = &Event{
		&calendar.Event{
			Summary: "eventBrokenDateTime",
			Start:   &calendar.EventDateTime{DateTime: t0.Format(time.RFC3339)},
			End:     &calendar.EventDateTime{DateTime: "---"},
		},
	}
)

func TestNewEvents(t *testing.T) {
	arg := []*calendar.Event{
		{
			Summary: "A",
			Start:   &calendar.EventDateTime{DateTime: t0.Format(time.RFC3339)},
			End:     &calendar.EventDateTime{DateTime: t0.Add(25 * time.Minute).Format(time.RFC3339)},
		},
		{
			Summary: "B",
			Start:   &calendar.EventDateTime{DateTime: t0.Add(60 * time.Minute).Format(time.RFC3339)},
			End:     &calendar.EventDateTime{DateTime: t0.Add(85 * time.Minute).Format(time.RFC3339)},
		},
		{
			Summary: "C",
			Start:   &calendar.EventDateTime{DateTime: t0.Add(120 * time.Minute).Format(time.RFC3339)},
			End:     &calendar.EventDateTime{DateTime: t0.Add(145 * time.Minute).Format(time.RFC3339)},
		},
	}
	want := Events{A, B, C}

	if diff := cmp.Diff(NewEvents(arg), want); diff != "" {
		t.Errorf("got an unexpected diff:\n%s", diff)
	}
}

func TestEvents_Diff(t *testing.T) {
	one := Events{A, B, C}
	another := Events{A, B, D}

	wantNegative := Events{D}
	wantPositive := Events{C}

	gotNegative, gotPositive, err := one.Diff(another)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

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
		expectError bool
	}{
		{
			name:   "True",
			events: Events{A},
			args: args{
				b: A,
			},
			want:        true,
			expectError: false,
		},
		{
			name:   "False",
			events: Events{A},
			args: args{
				b: B,
			},
			want:        false,
			expectError: false,
		},
		{
			name:   "BrokenStartDatetimeInEvents",
			events: Events{eventBrokenStartDateTime},
			args: args{
				b: A,
			},
			want:        false,
			expectError: true,
		},
		{
			name:   "BrokenEndDatetimeInEvents",
			events: Events{eventBrokenEndDateTime},
			args: args{
				b: A,
			},
			want:        false,
			expectError: true,
		},
		{
			name:   "BrokenStartDatetimeInB",
			events: Events{A},
			args: args{
				b: eventBrokenStartDateTime,
			},
			want:        false,
			expectError: true,
		},
		{
			name:   "BrokenEndDatetimeInB",
			events: Events{A},
			args: args{
				b: eventBrokenEndDateTime,
			},
			want:        false,
			expectError: true,
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			got, err := c.events.Has(c.args.b)

			if c.expectError {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if diff := cmp.Diff(got, c.want); diff != "" {
				t.Errorf("got an unexpected diff:\n%s", diff)
			}
		})
	}
}

func TestEvents_Unique(t *testing.T) {
	redundant := Events{A, A, B, B, C, C}
	want := Events{A, B, C}

	got, err := redundant.Unique()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("got an unexpected diff:\n%s", diff)
	}
}

func TestFetch(t *testing.T) {
	// TODO: change implementation of Fetch so that you can test easily
}

package event

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/api/calendar/v3"
)

func TestNew(t *testing.T) {
	type args struct {
		d     string
		s     string
		e     string
		title string
	}

	tests := []struct {
		name string
		args args
		want *Event
	}{
		{
			name: "Success",
			args: args{
				d:     "2022-09-20",
				s:     "07:00",
				e:     "07:25",
				title: "TEST_TITLE",
			},
			want: &Event{
				&calendar.Event{
					Summary: "TEST_TITLE",
					Start: &calendar.EventDateTime{
						DateTime: "2022-09-20T7:00:00+09:00",
					},
					End: &calendar.EventDateTime{
						DateTime: "2022-09-20T7:25:00+09:00",
					},
				},
			},
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			if diff := cmp.Diff(New(c.args.d, c.args.s, c.args.e, c.args.title), c.want); diff != "" {
				t.Errorf("got an unexpected diff:\n%s", diff)
			}
		})
	}
}

func TestEvent_Equals(t *testing.T) {
	tests := []struct {
		name        string
		event       *Event
		arg         *Event
		want        bool
		expectError bool
		expectPanic bool
	}{
		{
			name: "Equal",
			event: &Event{
				&calendar.Event{
					Summary: "TEST_TITLE",
					Start: &calendar.EventDateTime{
						DateTime: "2022-09-20T7:00:00+09:00",
					},
					End: &calendar.EventDateTime{
						DateTime: "2022-09-20T7:25:00+09:00",
					},
				},
			},
			arg: &Event{
				&calendar.Event{
					Summary: "TEST_TITLE",
					Start: &calendar.EventDateTime{
						DateTime: "2022-09-20T7:00:00+09:00",
					},
					End: &calendar.EventDateTime{
						DateTime: "2022-09-20T7:25:00+09:00",
					},
				},
			},
			want:        true,
			expectPanic: false,
		},
		{
			name: "DifferentSummary",
			event: &Event{
				&calendar.Event{
					Summary: "TEST_TITLE",
					Start: &calendar.EventDateTime{
						DateTime: "2022-09-20T7:00:00+09:00",
					},
					End: &calendar.EventDateTime{
						DateTime: "2022-09-20T7:25:00+09:00",
					},
				},
			},
			arg: &Event{
				&calendar.Event{
					Summary: "DIFFERENT_TITLE",
					Start: &calendar.EventDateTime{
						DateTime: "2022-09-20T7:00:00+09:00",
					},
					End: &calendar.EventDateTime{
						DateTime: "2022-09-20T7:25:00+09:00",
					},
				},
			},
			want:        false,
			expectPanic: false,
		},
		{
			name: "DifferentStartDateTime",
			event: &Event{
				&calendar.Event{
					Summary: "TEST_TITLE",
					Start: &calendar.EventDateTime{
						DateTime: "2022-09-20T7:00:00+09:00",
					},
					End: &calendar.EventDateTime{
						DateTime: "2022-09-20T7:25:00+09:00",
					},
				},
			},
			arg: &Event{
				&calendar.Event{
					Summary: "TEST_TITLE",
					Start: &calendar.EventDateTime{
						DateTime: "2022-09-19T7:00:00+09:00",
					},
					End: &calendar.EventDateTime{
						DateTime: "2022-09-20T7:25:00+09:00",
					},
				},
			},
			want:        false,
			expectPanic: false,
		},
		{
			name: "DifferentEndDateTime",
			event: &Event{
				&calendar.Event{
					Summary: "TEST_TITLE",
					Start: &calendar.EventDateTime{
						DateTime: "2022-09-20T7:00:00+09:00",
					},
					End: &calendar.EventDateTime{
						DateTime: "2022-09-20T7:25:00+09:00",
					},
				},
			},
			arg: &Event{
				&calendar.Event{
					Summary: "TEST_TITLE",
					Start: &calendar.EventDateTime{
						DateTime: "2022-09-20T7:00:00+09:00",
					},
					End: &calendar.EventDateTime{
						DateTime: "2022-09-21T7:25:00+09:00",
					},
				},
			},
			want:        false,
			expectPanic: false,
		},
		{
			name: "Panic",
			event: &Event{
				&calendar.Event{
					Summary: "TEST_TITLE",
					Start: &calendar.EventDateTime{
						DateTime: "2022-09-20T7:00:00+09:00",
					},
					End: &calendar.EventDateTime{
						DateTime: "2022-09-20T7:25:00+09:00",
					},
				},
			},
			arg: &Event{
				&calendar.Event{
					Summary: "TEST_TITLE",
					Start:   nil,
					End: &calendar.EventDateTime{
						DateTime: "2022-09-21T7:25:00+09:00",
					},
				},
			},
			want:        false,
			expectPanic: true,
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			defer func() {
				r := recover()

				switch {
				case r != nil && c.expectPanic:
					// OK
				case r != nil && !c.expectPanic:
					t.Errorf("unexpected panic: %v", r)
				case r == nil && c.expectPanic:
					t.Error("expected panic but did not panic")
				case r == nil && !c.expectPanic:
					// OK
				}
			}()

			got, err := c.event.Equals(c.arg)
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

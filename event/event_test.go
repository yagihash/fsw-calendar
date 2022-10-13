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
		want *calendar.Event
	}{
		{
			name: "Success",
			args: args{
				d:     "2022-09-20",
				s:     "07:00",
				e:     "07:25",
				title: "TEST_TITLE",
			},
			want: &calendar.Event{
				Summary: "TEST_TITLE",
				Start: &calendar.EventDateTime{
					DateTime: "2022-09-20T7:00:00+09:00",
				},
				End: &calendar.EventDateTime{
					DateTime: "2022-09-20T7:25:00+09:00",
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

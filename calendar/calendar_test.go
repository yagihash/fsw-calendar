package calendar

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNextMonth(t *testing.T) {
	cases := []struct {
		y  int
		m  int
		ny int
		nm int
	}{
		{2000, 1, 2000, 2},
		{2000, 12, 2001, 1},
	}

	for _, c := range cases {
		ny, nm := NextMonth(c.y, c.m)

		if diff := cmp.Diff(ny, c.ny); diff != "" {
			t.Errorf("unexpected next year\n%s", diff)
		}

		if diff := cmp.Diff(nm, c.nm); diff != "" {
			t.Errorf("unexpected next month\n%s", diff)
		}
	}
}

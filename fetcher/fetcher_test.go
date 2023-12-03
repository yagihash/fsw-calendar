package fetcher

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/yagihash/fsw-calendar/fetcher/class"
	"github.com/yagihash/fsw-calendar/fetcher/course"
)

func TestNew(t *testing.T) {
	New("example.com", course.RC, class.S4, http.DefaultClient)
}

func TestFetcher_FetchDocEvents(t *testing.T) {
	t.Run("success to fetch", func(t *testing.T) {
		c := NewTestClient(func(req *http.Request) *http.Response {
			switch req.URL.Path {
			case "/driving/sports/rc/t-4/2023/11.html":
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer([]byte("<table id=\"table-calendar\"><tbody><tr id=\"2023-11-2\" class=\"row-rc is-show\" data-date=\"2023-11-02\"><td class=\"day this-month thu\" data-date=\"2023/11/02\"><div class=\"inner-table-cell\">2<span class=\"yobi\">（木<span class=\"holiday-text\"></span>）</span></div></td><td class=\"type\"><div class=\"rc\"><p class=\"cell\">T-4 X</p></div></td><td class=\"time\"><div class=\"rc\"><p class=\"cell\">15:20~15:50</p></div></td><td class=\"wait\"><div class=\"rc\"><p class=\"cell\">30分</p></div></td></tr></tbody></table>"))),
					Header:     make(http.Header),
				}
			case "/driving/sports/rc/t-4/2023/12.html":
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer([]byte("<table id=\"table-calendar\"><tbody><tr id=\"2023-12-1\" class=\"row-rc is-show\" data-date=\"2023-12-01\"><td class=\"day this-month fri\" data-date=\"2023/12/01\"><div class=\"inner-table-cell\">2<span class=\"yobi\">（金<span class=\"holiday-text\"></span>）</span></div></td><td class=\"type\"><div class=\"rc\"><p class=\"cell\">T-4 X</p></div></td><td class=\"time\"><div class=\"rc\"><p class=\"cell\">15:20~15:50</p></div></td><td class=\"wait\"><div class=\"rc\"><p class=\"cell\">30分</p></div></td></tr></tbody></table>"))),
					Header:     make(http.Header),
				}
			default:
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer([]byte(""))),
					Header:     make(http.Header),
				}
			}
		})
		f := New("example.com", course.RC, class.T4, c)
		got, err := f.FetchDocEvents(2023, 11, 2)
		want := []DocEvent{{Date: "2023-11-02", Start: "15:20", End: "15:50", Title: "T-4 X"}, {Date: "2023-12-01", Start: "15:20", End: "15:50", Title: "T-4 X"}}

		if err != nil {
			t.Errorf("got an unexpexted error: %v", err)
		}

		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("got unexpected diff:\n%s", diff)
		}
	})

	t.Run("error http not found", func(t *testing.T) {
		c := NewTestClient(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(bytes.NewBuffer([]byte(""))),
				Header:     make(http.Header),
			}
		})
		f := New("example.com", course.RC, class.T4, c)
		_, err := f.FetchDocEvents(2023, 11, 1)

		if err.Error() != "got status code 404 on https://example.com/driving/sports/rc/t-4/2023/11.html" {
			t.Errorf("got an unexpected error: %v", err)
		}
	})
}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

package frame

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Frames []*Frame

func (fs Frames) Diff(another Frames) (negative, positive Frames) {
	if reflect.DeepEqual(fs, another) {
		return
	}

	for _, f := range another {
		if !fs.Has(f) {
			negative = append(negative, f)
		}
	}

	for _, f := range fs {
		if !another.Has(f) {
			positive = append(positive, f)
		}
	}

	return
}

func (fs Frames) Has(b *Frame) bool {
	for _, a := range fs {
		if reflect.DeepEqual(a, b) {
			return true
		}
	}

	return false
}

func (fs Frames) Add(f *Frame) {
	fs = append(fs, f)
}

func Fetch(url string) (Frames, error) {
	var frames []*Frame

	res, err := http.Get(url)
	if err != nil {
		return frames, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return frames, fmt.Errorf("status code error: %s", res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return frames, err
	}

	doc.Find("#table-calendar > tbody > tr.row-rc > td.type > div > p").Each(func(i int, s *goquery.Selection) {
		d, _ := s.Parent().Parent().Parent().Attr("data-date")
		t := strings.Split(doc.Find("#table-calendar > tbody > tr.row-rc > td.time > div > p").Eq(i).Text(), "~")
		frames = append(frames, New(d, t[0], t[1], s.Text()))
	})

	return frames, nil
}

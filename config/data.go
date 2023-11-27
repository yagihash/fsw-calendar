package config

import (
	"encoding/json"

	"github.com/yagihash/fsw-calendar/fetcher/class"
	"github.com/yagihash/fsw-calendar/fetcher/course"
)

type Data struct {
	CalendarID string        `json:"calendar_id"`
	Course     course.Course `json:"course"`
	Class      class.Class   `json:"class""`
}

func (d *Data) UnmarshalJSON(b []byte) error {
	var tmp map[string]string

	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	switch tmp["course"] {
	case "rc":
		d.Course = course.RC
	case "ss":
		d.Course = course.SS
	default:
		d.Course = course.Unknown
	}

	switch tmp["class"] {
	case "ss-4":
		d.Class = class.SS4
	case "t-4":
		d.Class = class.T4
	case "ns-4":
		d.Class = class.NS4
	case "s-4":
		d.Class = class.S4
	default:
		d.Class = class.Unknown
	}

	return nil
}

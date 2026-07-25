package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Date struct {
	time.Time
}

const dateLayout = "2006-01-02"

func (d *Date) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		return nil
	}
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return fmt.Errorf("invalid date %q, expected format %s: %w",
			s, dateLayout, err)
	}
	d.Time = t
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.Format(dateLayout))
}

func ToTimePtr(d *Date) *time.Time {
	if d == nil {
		return nil
	}
	t := d.Time
	return &t
}

package types

import (
	"encoding/json"
	"time"
)

const layout = "2006-01-02"

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(b []byte) error {

	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	if s == "" {
		d.Time = time.Time{}
		return nil
	}

	t, err := time.Parse(layout, s)
	if err != nil {
		return err
	}

	d.Time = t
	return nil
}

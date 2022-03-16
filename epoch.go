package main

import (
	"encoding/json"
	"time"
)

type timestamp struct {
	time.Time
}

func (p *timestamp) UnmarshalJSON(bytes []byte) error {
	var raw int64
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	p.Time = time.Unix(raw, 0)
	return nil
}

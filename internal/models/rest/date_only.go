package rest

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type DateOnly time.Time

const DateOnlyFormat = "2006-01-02"

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), "\"")
	t, err := time.Parse(DateOnlyFormat, str)
	if err != nil {
		return err
	}
	*d = DateOnly(t)
	return nil
}

func (dt DateOnly) MarshalJSON() ([]byte, error) {
	time := time.Time(dt)
	year, m, d := time.Year(), time.Month(), time.Day()
	var month string
	if m > 10 {
		month = fmt.Sprintf("%d", m)
	} else {
		month = fmt.Sprintf("0%d", m)
	}
	var day string
	if d > 10 {
		day = fmt.Sprintf("%d", d)
	} else {
		day = fmt.Sprintf("0%d", d)
	}
	return json.Marshal(fmt.Sprintf("%d-%s-%s", year, month, day))
}

func (dt DateOnly) Format(s string) string {
	return time.Time(dt).Format(s)
}

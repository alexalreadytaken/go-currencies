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

func (d DateOnly) MarshalJSON() ([]byte, error) {
	year, month, day := time.Time(d).Date()
	return json.Marshal(fmt.Sprintf("%d-%d-%d", year, month, day))
}

func (d DateOnly) Format(s string) string {
	return time.Time(d).Format(s)
}

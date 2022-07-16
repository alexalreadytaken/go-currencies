package rest

import "encoding/json"

type BtcUsdtCourseSlice struct {
	Value     float64 `json:"value"`
	Timestamp uint64  `json:"timestamp"`
}

type BtcFiatCourseSlice map[string]float64

// @description has additional fields consist of currency code and amount, but cant image here
type FiatToAnyCourseSlice struct {
	Currencies   map[string]float64 `json:"-"`
	BaseCurrency string             `json:"base_currency"`
	Date         DateOnly           `json:"date"`
}

func (slice FiatToAnyCourseSlice) MarshalJSON() ([]byte, error) {
	type slice_ FiatToAnyCourseSlice
	b, err := json.Marshal(slice_(slice))
	if err != nil {
		return nil, err
	}
	var m map[string]json.RawMessage
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	for k, v := range slice.Currencies {
		b, err = json.Marshal(v)
		if err != nil {
			return nil, err
		}
		m[k] = b
	}
	return json.Marshal(m)
}

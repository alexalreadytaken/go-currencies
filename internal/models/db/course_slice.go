package db

import (
	"time"
)

type BtcUsdtCourseSlice struct {
	ID        uint
	Value     float64
	Timestamp time.Time
}

type AnyToFiatCourseSlice struct {
	ID               uint
	BaseCurrencyCode string
	FiatCode         string
	Value            float64
	GetDate          time.Time
}

func (slice *BtcUsdtCourseSlice) TableName() string {
	return "btc_usdt_course"
}

func (slice *AnyToFiatCourseSlice) TableName() string {
	return "any_to_fiat_course"
}

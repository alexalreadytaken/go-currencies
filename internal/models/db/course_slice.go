package db

import (
	"time"

	"gorm.io/gorm"
)

type BtcUsdtCourseSlice struct {
	gorm.Model
	Value     float64
	Timestamp time.Time
}

type AnyToFiatCourseSlice struct {
	gorm.Model
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

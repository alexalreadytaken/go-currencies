package utils

import (
	"time"

	"github.com/alexalreadytaken/go-currencies/internal/models/db"
	"github.com/alexalreadytaken/go-currencies/internal/models/rest"
	"github.com/stretchr/testify/mock"
)

type MockCurrenciesCourseProducer struct {
	mock.Mock
}

func (producer *MockCurrenciesCourseProducer) GetBtcUsdCourse() (float64, error) {
	res := producer.Called()
	return res.Get(0).(float64), res.Error(0)
}

func (producer *MockCurrenciesCourseProducer) GetRubToFiatCourse() (*rest.FiatPricesResposnse, error) {
	res := producer.Called()
	return res.Get(0).(*rest.FiatPricesResposnse), res.Error(0)
}

type MockCurrenciesRepo struct {
	mock.Mock
}

func (repo *MockCurrenciesRepo) GetLastBtcUsdt() (*db.BtcUsdtCourseSlice, error) {
	res := repo.Called()
	return res.Get(0).(*db.BtcUsdtCourseSlice), res.Error(1)
}

func (repo *MockCurrenciesRepo) GetBtcUsdtHistory(
	limit uint, offset uint, fromTime time.Time) (
	int64, []db.BtcUsdtCourseSlice, error) {
	res := repo.Called(limit, offset, fromTime)
	return res.Get(0).(int64), res.Get(1).([]db.BtcUsdtCourseSlice), res.Error(2)
}

func (repo *MockCurrenciesRepo) GetLastAnyToFiat(baseCurrencyCode string) ([]db.AnyToFiatCourseSlice, error) {
	res := repo.Called(baseCurrencyCode)
	return res.Get(0).([]db.AnyToFiatCourseSlice), res.Error(1)
}

func (repo *MockCurrenciesRepo) GetAnyToFiatHistory(
	baseCurrencyCode string, limit uint, offset uint, fromDate time.Time) (
	int64, []db.AnyToFiatCourseSlice, error) {
	res := repo.Called(baseCurrencyCode, limit, offset, fromDate)
	return res.Get(0).(int64), res.Get(1).([]db.AnyToFiatCourseSlice), res.Error(2)
}

func (repo *MockCurrenciesRepo) GetLastSavedDateAnyToFiatCourse(baseCurrencyCode string) (*time.Time, error) {
	res := repo.Called(baseCurrencyCode)
	return res.Get(0).(*time.Time), res.Error(1)
}

func (repo *MockCurrenciesRepo) StoreBtcUsdtSlice(slice db.BtcUsdtCourseSlice) error {
	res := repo.Called(slice)
	return res.Error(0)
}

func (repo *MockCurrenciesRepo) StoreAnyToFiatSlice(slices []db.AnyToFiatCourseSlice) error {
	res := repo.Called(slices)
	return res.Error(0)
}
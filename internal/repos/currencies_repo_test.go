package repos

import (
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alexalreadytaken/go-currencies/internal/models/db"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AnyTime struct{}

func (t AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type GormCurrenciesRepoTestSuite struct {
	suite.Suite
	repo *GormCurrenciesRepo
	mock sqlmock.Sqlmock
}

func TestGormCurrenciesRepo(t *testing.T) {
	suite.Run(t, &GormCurrenciesRepoTestSuite{})
}

func (s *GormCurrenciesRepoTestSuite) SetupSuite() {
	db, mock, err := sqlmock.New()
	s.NoError(err)

	pg := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gormDb, err := gorm.Open(pg, &gorm.Config{})
	s.NoError(err)
	s.mock = mock
	s.repo = &GormCurrenciesRepo{
		db: gormDb,
	}
}

func (s *GormCurrenciesRepoTestSuite) TestGetLastBtcUsdt() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "btc_usdt_course"
			ORDER BY timestamp desc LIMIT 1
		`)).WillReturnRows(sqlmock.NewRows([]string{"id", "value", "timestamp"}))
	s.repo.GetLastBtcUsdt()
}

func (s *GormCurrenciesRepoTestSuite) TestGetAnyToFiatHistory() {
	code, limit, offset, fromTime := "BTC", uint(0), uint(0), time.Now()
	// fixme
	// LIMIT $3 OFFSET $4
	s.mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT DISTINCT "get_date" FROM "any_to_fiat_course"
			WHERE base_currency_code = $1 AND get_date > $2
			ORDER BY get_date
		`)).
		WithArgs(code, fromTime).
		WillReturnRows(sqlmock.NewRows([]string{"get_date"}))
	//fixme (NULL)
	s.mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "any_to_fiat_course"
			WHERE get_date in (NULL)
			ORDER BY get_date`)).
		WillReturnRows(sqlmock.
			NewRows([]string{"base_currency_code", "fiat_code", "value", "get_date"}))
	s.repo.GetAnyToFiatHistory(code, limit, offset, fromTime)
}

func (s *GormCurrenciesRepoTestSuite) TestGetLastSavedDateAnyToFiatCourse() {
	code := "BTC"
	s.mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT DISTINCT "get_date" FROM "any_to_fiat_course"
		WHERE base_currency_code = $1
		ORDER BY get_date desc LIMIT 1
	`)).WithArgs(code).
		WillReturnRows(sqlmock.NewRows([]string{"get_date"}))
	s.repo.GetLastSavedDateAnyToFiatCourse(code)
}

func (s *GormCurrenciesRepoTestSuite) TestGetBtcUsdtHistory() {
	limit, offset, fromTime := uint(0), uint(0), time.Now()
	s.mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT count(*) FROM "btc_usdt_course" WHERE timestamp >= $1
		`)).WithArgs(fromTime).WillReturnRows(sqlmock.NewRows([]string{"count"}))
	//fixme
	// desc LIMIT $2 OFFSET $3
	s.mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "btc_usdt_course" WHERE timestamp >= $1
			ORDER BY timestamp`)).
		WithArgs(fromTime).
		WillReturnRows(sqlmock.NewRows([]string{"value", "timestamp"}))
	s.repo.GetBtcUsdtHistory(limit, offset, fromTime)
}

func (s *GormCurrenciesRepoTestSuite) TestGetLastAnyToFiat() {
	code := "BTC"
	s.mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT "get_date" FROM "any_to_fiat_course"
			WHERE base_currency_code = $1
			ORDER BY get_date desc LIMIT 1`)).
		WithArgs(code).
		WillReturnRows(sqlmock.NewRows([]string{"get_date"}))
	s.mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "any_to_fiat_course"
			WHERE base_currency_code = $1 AND get_date = $2`)).
		WithArgs(code, AnyTime{}).
		WillReturnRows(sqlmock.NewRows(
			[]string{"base_currency_code", "fiat_code", "value", "get_date"}))
	s.repo.GetLastAnyToFiat(code)
}

func (s *GormCurrenciesRepoTestSuite) TestStoreAnyToFiatSlice() {
	slice := db.AnyToFiatCourseSlice{
		BaseCurrencyCode: "RUB",
		FiatCode:         "USD",
		Value:            100,
		GetDate:          time.Now(),
	}
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`
			INSERT INTO "any_to_fiat_course" ("base_currency_code","fiat_code","value","get_date")
			VALUES ($1,$2,$3,$4) RETURNING "id"`)).
		WithArgs(slice.BaseCurrencyCode, slice.FiatCode, slice.Value, slice.GetDate).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	s.mock.ExpectCommit()
	s.repo.StoreAnyToFiatSlice([]db.AnyToFiatCourseSlice{slice})
}

func (s *GormCurrenciesRepoTestSuite) TestStoreBtcUsdtSlice() {
	slice := db.BtcUsdtCourseSlice{
		Value:     100.10,
		Timestamp: time.Now(),
	}
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`
			INSERT INTO "btc_usdt_course" ("value","timestamp")
			VALUES ($1,$2) RETURNING "id"`)).
		WithArgs(slice.Value, slice.Timestamp).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	s.mock.ExpectCommit()
	s.repo.StoreBtcUsdtSlice(slice)
}

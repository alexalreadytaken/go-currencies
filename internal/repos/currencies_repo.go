package repos

import (
	"fmt"
	"time"

	"github.com/alexalreadytaken/go-currencies/internal/models/db"
	"github.com/alexalreadytaken/go-currencies/internal/utils"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CurrenciesRepo struct {
	db *gorm.DB
}

var (
	btcusdModel    = &db.BtcUsdtCourseSlice{}
	anyToFiatModel = &db.AnyToFiatCourseSlice{}
)

func NewCurrenciesRepo(cnf *utils.AppConfig) (*CurrenciesRepo, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cnf.DbHost, cnf.DbPort, cnf.DbUser, cnf.DbPass, cnf.DbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(btcusdModel, anyToFiatModel)
	if err != nil {
		return nil, err
	}
	return &CurrenciesRepo{
		db: db,
	}, nil
}

func (repo *CurrenciesRepo) GetLastBtcUsdt() (*db.BtcUsdtCourseSlice, error) {
	var courseSlice db.BtcUsdtCourseSlice
	err := repo.db.Model(btcusdModel).
		Order("timestamp desc").
		Limit(1).
		Find(&courseSlice).
		Error
	if err != nil {
		return nil, err
	}
	return &courseSlice, nil
}

func (repo *CurrenciesRepo) GetBtcUsdtHistory(
	limit uint, offset uint, fromTime time.Time) (
	total int64, history []db.BtcUsdtCourseSlice, err error) {
	err = repo.db.Model(btcusdModel).
		Where("timestamp >= ?", fromTime).
		Count(&total).
		Error
	if err != nil {
		return 0, nil, err
	}
	err = repo.db.Model(btcusdModel).
		Order("timestamp").
		Where("timestamp >= ?", fromTime).
		Limit(int(limit)).
		Offset(int(offset)).
		Find(&history).
		Error
	if err != nil {
		return 0, nil, err
	}
	return
}

func (repo *CurrenciesRepo) GetLastAnyToFiat(baseCurrencyCode string) ([]db.AnyToFiatCourseSlice, error) {
	var courses []db.AnyToFiatCourseSlice
	var lastDate time.Time
	err := repo.db.Model(anyToFiatModel).
		Where("base_currency_code = ?", baseCurrencyCode).
		Select("get_date").
		Order("get_date desc").
		Limit(1).
		Find(&lastDate).Error
	if err != nil {
		return nil, err
	}
	err = repo.db.Model(anyToFiatModel).
		Where("base_currency_code = ?", baseCurrencyCode).
		Where("get_date = ?", lastDate).
		Find(&courses).
		Error
	if err != nil {
		return nil, err
	}
	return courses, nil
}

func (repo *CurrenciesRepo) GetAnyToFiatHistory(
	baseCurrencyCode string,
	limit uint,
	offset uint,
	fromDate datatypes.Date) (total int64, history []db.AnyToFiatCourseSlice, err error) {
	var dates []time.Time
	err = repo.db.Model(anyToFiatModel).
		Where("base_currency_code = ?", baseCurrencyCode).
		Where("get_date > ?", fromDate).
		Order("get_date").
		Distinct("get_date").
		Limit(int(limit)).
		Offset(int(offset)).
		Find(&dates).
		Error
	if err != nil {
		return 0, nil, err
	}
	total = int64(len(dates))
	err = repo.db.Model(anyToFiatModel).
		Where("get_date in ?", dates).
		Order("get_date").
		Find(&history).
		Error
	if err != nil {
		return 0, nil, err
	}
	return
}

func (repo *CurrenciesRepo) GetLastSavedDateAnyToFiatCourse(baseCurrencyCode string) (*time.Time, error) {
	var lastDate time.Time
	err := repo.db.Model(anyToFiatModel).
		Where("base_currency_code = ?", baseCurrencyCode).
		Order("get_date desc").
		Distinct("get_date").
		Select("get_date").
		Limit(1).
		Find(&lastDate).
		Error
	if err != nil {
		return nil, err
	}
	return &lastDate, nil
}

func (repo *CurrenciesRepo) StoreBtcUsdtSlice(slice db.BtcUsdtCourseSlice) error {
	return repo.db.Create(&slice).Error
}

func (repo *CurrenciesRepo) StoreAnyToFiatSlice(slices []db.AnyToFiatCourseSlice) error {
	return repo.db.Create(&slices).Error
}

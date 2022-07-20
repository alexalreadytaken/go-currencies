package repos

import (
	"fmt"
	"time"

	"github.com/alexalreadytaken/go-currencies/internal/models/db"
	"github.com/alexalreadytaken/go-currencies/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CurrenciesRepo interface {
	GetLastBtcUsdt() (*db.BtcUsdtCourseSlice, error)

	GetBtcUsdtHistory(limit uint, offset uint, fromTime time.Time) (int64, []db.BtcUsdtCourseSlice, error)

	GetLastAnyToFiat(baseCurrencyCode string) ([]db.AnyToFiatCourseSlice, error)

	GetAnyToFiatHistory(baseCurrencyCode string, limit uint, offset uint, fromDate time.Time) (int64, []db.AnyToFiatCourseSlice, error)

	GetLastSavedDateAnyToFiatCourse(baseCurrencyCode string) (*time.Time, error)

	StoreBtcUsdtSlice(slice db.BtcUsdtCourseSlice) error

	StoreAnyToFiatSlice(slices []db.AnyToFiatCourseSlice) error
}

type GormCurrenciesRepo struct {
	db *gorm.DB
}

var (
	btcusdModel    = &db.BtcUsdtCourseSlice{}
	anyToFiatModel = &db.AnyToFiatCourseSlice{}
)

func NewCurrenciesRepo(cnf *utils.AppConfig) (*GormCurrenciesRepo, error) {
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
	return &GormCurrenciesRepo{
		db: db,
	}, nil
}

func (repo *GormCurrenciesRepo) GetLastBtcUsdt() (*db.BtcUsdtCourseSlice, error) {
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

func (repo *GormCurrenciesRepo) GetBtcUsdtHistory(
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
		Limit(int(limit)).
		Offset(int(offset)).
		Order("timestamp").
		Where("timestamp >= ?", fromTime).
		Find(&history).
		Error
	if err != nil {
		return 0, nil, err
	}
	return
}

func (repo *GormCurrenciesRepo) GetLastAnyToFiat(baseCurrencyCode string) ([]db.AnyToFiatCourseSlice, error) {
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

func (repo *GormCurrenciesRepo) GetAnyToFiatHistory(
	baseCurrencyCode string,
	limit uint,
	offset uint,
	fromDate time.Time) (total int64, history []db.AnyToFiatCourseSlice, err error) {
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

func (repo *GormCurrenciesRepo) GetLastSavedDateAnyToFiatCourse(baseCurrencyCode string) (*time.Time, error) {
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

func (repo *GormCurrenciesRepo) StoreBtcUsdtSlice(slice db.BtcUsdtCourseSlice) error {
	return repo.db.Create(&slice).Error
}

func (repo *GormCurrenciesRepo) StoreAnyToFiatSlice(slices []db.AnyToFiatCourseSlice) error {
	return repo.db.Create(&slices).Error
}

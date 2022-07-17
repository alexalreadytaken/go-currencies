package services

import (
	"fmt"
	"log"
	"time"

	"github.com/alexalreadytaken/go-currencies/internal/models/db"
	"github.com/alexalreadytaken/go-currencies/internal/repos"
	"github.com/alexalreadytaken/go-currencies/internal/utils"
)

type CurrencyDaemonWorker struct {
	client *CurrenciesClient
	repo   *repos.CurrenciesRepo

	latestBtcUsdCourse float64
}

func NewCurrencyDaemonWorker(client *CurrenciesClient, repo *repos.CurrenciesRepo) (*CurrencyDaemonWorker, error) {
	latest, err := repo.GetLastBtcUsdt()
	if err != nil {
		return nil, fmt.Errorf("can't get info about latest btc usd course for caching = %s", err.Error())
	}
	return &CurrencyDaemonWorker{
		latestBtcUsdCourse: latest.Value,
		repo:               repo,
		client:             client,
	}, nil
}

func (worker *CurrencyDaemonWorker) GetBtcUsdtCourseAndStore() {
	course, err := worker.client.GetBtcUsdCourse()
	t := time.Now()
	if err != nil {
		log.Printf("error while getting btc usd course = %s", err.Error())
		return
	}
	if course == worker.latestBtcUsdCourse {
		return
	}
	err = worker.repo.StoreBtcUsdtSlice(db.BtcUsdtCourseSlice{
		Value:     course,
		Timestamp: t,
	})
	if err != nil {
		log.Printf("error while saving btc usd course = %s", err.Error())
		return
	}
	worker.updateBtcToFiatCourseAndStore(course)
	worker.latestBtcUsdCourse = course
}

func (worker *CurrencyDaemonWorker) updateBtcToFiatCourseAndStore(newBtcUsdtCourse float64) {
	slices, err := worker.repo.GetLastAnyToFiat("RUB")
	if err != nil {
		log.Printf("error while getting info about rub to fiat course = %s", err.Error())
		return
	}
	if len(slices) == 0 {
		log.Println("can`t find info about rub to fiat crourses")
		return
	}
	var rubUsdCourse float64
	for i := 0; i < len(slices); i++ {
		slice := slices[i]
		if slice.FiatCode == "USD" {
			rubUsdCourse = slice.Value
		}
	}
	t := time.Now()
	var btcToFiatSlices []db.AnyToFiatCourseSlice
	for i := 0; i < len(slices); i++ {
		rubSlice := slices[i]
		btcSlice := db.AnyToFiatCourseSlice{
			BaseCurrencyCode: "BTC",
			FiatCode:         rubSlice.FiatCode,
			GetDate:          t,
			Value:            (((1 / rubUsdCourse) * newBtcUsdtCourse) * rubSlice.Value),
		}
		btcToFiatSlices = append(btcToFiatSlices, btcSlice)
	}
	btcToFiatSlices = append(btcToFiatSlices, db.AnyToFiatCourseSlice{
		BaseCurrencyCode: "BTC",
		FiatCode:         "RUB",
		GetDate:          t,
		Value:            (1 / rubUsdCourse) * newBtcUsdtCourse,
	})
	err = worker.repo.StoreAnyToFiatSlice(btcToFiatSlices)
	if err != nil {
		log.Printf("error while saving btc to fiat course = %s", err.Error())
	}
}

func (worker *CurrencyDaemonWorker) GetRubToFiatCourseAndStore() {
	t := time.Now()
	lastDate, err := worker.repo.GetLastSavedDateAnyToFiatCourse("RUB")
	if err != nil {
		log.Printf("error while getting last saved date rub to fiat course = %s", err.Error())
		return
	}
	if utils.DatesEquals(*lastDate, t) {
		log.Println("rub to fiat course today already saved, skip")
		return
	}
	course, err := worker.client.GetRubToFiatCourse()
	if err != nil {
		log.Printf("error while getting rub to fiat course = %s", err.Error())
		return
	}
	var slices []db.AnyToFiatCourseSlice
	for k, v := range course.Rates {
		slices = append(slices, db.AnyToFiatCourseSlice{
			BaseCurrencyCode: course.Base,
			FiatCode:         k,
			Value:            v,
			GetDate:          t,
		})
	}
	err = worker.repo.StoreAnyToFiatSlice(slices)
	if err != nil {
		log.Printf("error while saving rub to fiat courses = %s", err.Error())
	}
}

package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexalreadytaken/go-currencies/internal/models/db"
	"github.com/alexalreadytaken/go-currencies/internal/models/rest"
	"github.com/alexalreadytaken/go-currencies/internal/repos"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type CurrenciesController struct {
	repo *repos.CurrenciesRepo
}

func NewCurrenciesController(repo *repos.CurrenciesRepo) *CurrenciesController {
	return &CurrenciesController{
		repo: repo,
	}
}

func (c *CurrenciesController) GetLastBtcUsdtCourse(g *gin.Context) {
	lastCourse, err := c.repo.GetLastBtcUsdt()
	if err != nil {
		msg := "error while getting last btc usdt course"
		log.Println(msg, err)
		g.AbortWithStatusJSON(http.StatusInternalServerError,
			rest.UnexpectedResult{Message: msg})
		return
	}
	g.JSON(http.StatusOK, rest.BtcUsdtCourseSlice{
		Value:     lastCourse.Value,
		Timestamp: uint64(lastCourse.Timestamp.Unix()),
	})
}

func (c *CurrenciesController) GetBtcUsdtCourseHistory(g *gin.Context) {
	var pagination rest.BtcUsdtPaginationRequest
	if err := g.Bind(&pagination); err != nil {
		msg := "invalid pagination format="
		log.Println(msg, err)
		g.AbortWithStatusJSON(http.StatusBadRequest,
			rest.UnexpectedResult{Message: msg + err.Error()})
		return
	}
	total, history, err := c.repo.GetBtcUsdtHistory(
		pagination.Limit,
		pagination.Offset,
		time.Unix(int64(pagination.FromTime), 0))
	if err != nil {
		msg := "error while getting info about btc usdt course"
		log.Println(msg, err)
		g.AbortWithStatusJSON(http.StatusInternalServerError,
			rest.UnexpectedResult{Message: msg})
		return
	}
	historyRest := make([]rest.BtcUsdtCourseSlice, len(history))
	for i := 0; i < len(history); i++ {
		sliceDb := history[i]
		historyRest[i] = rest.BtcUsdtCourseSlice{
			Value:     sliceDb.Value,
			Timestamp: uint64(sliceDb.Timestamp.Unix()),
		}
	}
	g.JSON(http.StatusOK, rest.BtcUsdtHistoryPage{
		Total:   int(total),
		History: historyRest,
	})
}

func (c *CurrenciesController) GetLastRubToFiatCourse(g *gin.Context) {
	code := "RUB"
	c.getLastAnyToFiat(g, code)
}

func (c *CurrenciesController) GetRubToFiatCourseHistory(g *gin.Context) {
	code := "RUB"
	c.getAnyToFiatHistory(g, code)
}

func (c *CurrenciesController) GetLastBtcToFiatCourse(g *gin.Context) {
	code := "BTC"
	c.getLastAnyToFiat(g, code)
}

func (c *CurrenciesController) getLastAnyToFiat(g *gin.Context, baseCurrencryCode string) {
	lastCourses, err := c.repo.GetLastAnyToFiat(baseCurrencryCode)
	if err != nil {
		msg := fmt.Sprintf("error while getting last %s to fiat course", baseCurrencryCode)
		log.Println(msg, err)
		g.AbortWithStatusJSON(http.StatusInternalServerError,
			rest.UnexpectedResult{Message: msg})
		return
	}
	if len(lastCourses) == 0 {
		g.AbortWithStatusJSON(http.StatusInternalServerError,
			rest.UnexpectedResult{Message: "not found"})
		return
	}
	currencies := make(map[string]float64, len(lastCourses))
	for i := 0; i < len(lastCourses); i++ {
		course := lastCourses[i]
		currencies[course.FiatCode] = course.Value
	}
	g.JSON(http.StatusOK, rest.AnyToFiatCourseSlice{
		Currencies:   currencies,
		BaseCurrency: baseCurrencryCode,
		Date:         rest.DateOnly(lastCourses[0].GetDate),
	})
}

func (c *CurrenciesController) getAnyToFiatHistory(g *gin.Context, baseCurrencryCode string) {
	var pagination rest.AnyToFiatPaginationRequest
	if err := g.Bind(&pagination); err != nil {
		msg := "invalid pagination format="
		log.Println(msg, err)
		g.AbortWithStatusJSON(http.StatusBadRequest,
			rest.UnexpectedResult{Message: msg + err.Error()})
		return
	}
	total, slices, err := c.repo.GetAnyToFiatHistory(
		baseCurrencryCode,
		pagination.Limit,
		pagination.Offset,
		datatypes.Date(pagination.FromDate))
	if err != nil {
		msg := "error while getting rub to fiat course"
		log.Println(msg, err)
		g.AbortWithStatusJSON(http.StatusInternalServerError,
			rest.UnexpectedResult{Message: msg})
		return
	}
	if len(slices) == 0 {
		g.AbortWithStatusJSON(http.StatusInternalServerError,
			rest.UnexpectedResult{Message: "not found"})
		return
	}
	//todo refactor
	dateToSlices := make(map[time.Time][]db.AnyToFiatCourseSlice)
	for i := 0; i < len(slices); i++ {
		slice := slices[i]
		slices := dateToSlices[slice.GetDate]
		slices = append(slices, slice)
		dateToSlices[slice.GetDate] = slices
	}
	var finalSlices []rest.AnyToFiatCourseSlice
	for date, slices := range dateToSlices {
		currencies := make(map[string]float64)
		for i := 0; i < len(slices); i++ {
			slice := slices[i]
			currencies[slice.FiatCode] = slice.Value
		}
		finalSlices = append(finalSlices, rest.AnyToFiatCourseSlice{
			BaseCurrency: baseCurrencryCode,
			Date:         rest.DateOnly(date),
			Currencies:   currencies,
		})
	}
	g.JSON(http.StatusOK, rest.AnyToFiatHistoryPage{
		Total:   int(total),
		History: finalSlices,
	})
}

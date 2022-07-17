package controllers

import "github.com/gin-gonic/gin"

func AddCurrenciesRoutes(api *gin.RouterGroup, controller *CurrenciesController) {
	btcusd := api.Group("/btcusdt")
	{
		btcusd.GET("", GetLastBtcUsd(controller))
		btcusd.POST("", GetBtcUsdHistory(controller))
	}
	currencies := api.Group("/currencies")
	{
		currencies.GET("", GetLastFiatRub(controller))
		currencies.POST("", GetFiatRubHistory(controller))
	}
	api.GET("/latest", GetLastBtcFiat(controller))
}

// Currencies godoc
// @Summary get last BTC USDT course
// @Schemes
// @Produce json
// @Success 200 {object} rest.BtcUsdtCourseSlice
// @Failure 400 {object} rest.UnexpectedResult
// @Failure 500 {object} rest.UnexpectedResult
// @Router /btcusdt [get]
func GetLastBtcUsd(controller *CurrenciesController) func(*gin.Context) {
	return controller.GetLastBtcUsdtCourse
}

// Currencies godoc
// @Summary get BTC USDT course history with pagination
// @Schemes
// @Produce json
// @Accept json
// @Param pagination body rest.BtcUsdtPaginationRequest true "pagination"
// @Success 200 {object} rest.BtcUsdtHistoryPage
// @Failure 400 {object} rest.UnexpectedResult
// @Failure 500 {object} rest.UnexpectedResult
// @Router /btcusdt [post]
func GetBtcUsdHistory(controller *CurrenciesController) func(*gin.Context) {
	return controller.GetBtcUsdtCourseHistory
}

// Currencies godoc
// @Summary get last RUB to fiat course
// @Schemes
// @Produce json
// @Success 200 {object} rest.AnyToFiatCourseSlice
// @Failure 400 {object} rest.UnexpectedResult
// @Failure 500 {object} rest.UnexpectedResult
// @Router /currencies [get]
func GetLastFiatRub(controller *CurrenciesController) func(*gin.Context) {
	return controller.GetLastRubToFiatCourse
}

// Currencies godoc
// @Summary get RUB to fiat course history with pagination
// @Schemes
// @Produce json
// @Accept json
// @Param pagination body rest.AnyToFiatPaginationRequest true "pagination"
// @Success 200 {object} rest.AnyToFiatHistoryPage
// @Failure 400 {object} rest.UnexpectedResult
// @Failure 500 {object} rest.UnexpectedResult
// @Router /currencies [post]
func GetFiatRubHistory(controller *CurrenciesController) func(*gin.Context) {
	return controller.GetRubToFiatCourseHistory
}

// Currencies godoc
// @Summary get last fiat to BTC course
// @Schemes
// @Produce json
// @Success 200 {object} rest.AnyToFiatCourseSlice
// @Failure 400 {object} rest.UnexpectedResult
// @Failure 500 {object} rest.UnexpectedResult
// @Router /latest [get]
func GetLastBtcFiat(controller *CurrenciesController) func(*gin.Context) {
	return controller.GetLastBtcToFiatCourse
}

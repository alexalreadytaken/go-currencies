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
// @Summary get last fiat to rub course
// @Schemes
// @Produce json
// @Success 200 {object} rest.FiatToAnyCourseSlice
// @Failure 400 {object} rest.UnexpectedResult
// @Failure 500 {object} rest.UnexpectedResult
// @Router /currencies [get]
func GetLastFiatRub(controller *CurrenciesController) func(*gin.Context) {
	return controller.GetLastBtcFiatCourse
}

// Currencies godoc
// @Summary get fiat to rub course history with pagination
// @Schemes
// @Produce json
// @Accept json
// @Param pagination body rest.FiatPaginationRequest true "pagination" 
// @Success 200 {object} rest.FiatToAnyHistoryPage
// @Failure 400 {object} rest.UnexpectedResult
// @Failure 500 {object} rest.UnexpectedResult
// @Router /currencies [post]
func GetFiatRubHistory(controller *CurrenciesController) func(*gin.Context) {
	return controller.GetFiatRubCourseHistory
}

// Currencies godoc
// @Summary get last fiat to BTC course
// @Schemes
// @Produce json
// @Success 200 {object} rest.FiatToAnyCourseSlice
// @Failure 400 {object} rest.UnexpectedResult
// @Failure 500 {object} rest.UnexpectedResult
// @Router /latest [get]
func GetLastBtcFiat(controller *CurrenciesController) func(*gin.Context) {
	return controller.GetLastBtcFiatCourse
}

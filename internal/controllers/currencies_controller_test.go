package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexalreadytaken/go-currencies/internal/models/db"
	"github.com/alexalreadytaken/go-currencies/internal/models/rest"
	"github.com/alexalreadytaken/go-currencies/internal/utils"
	ginpkg "github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type CurrenciesControllerTestSuite struct {
	suite.Suite
	gin      *ginpkg.Engine
	repoMock *utils.MockCurrenciesRepo
}

func TestCurrenciesController(t *testing.T) {
	suite.Run(t, &CurrenciesControllerTestSuite{})
}

func (s *CurrenciesControllerTestSuite) SetupSuite() {
	gin := ginpkg.Default()
	repo := &utils.MockCurrenciesRepo{}
	controller := NewCurrenciesController(repo)
	api := gin.Group("/api")
	AddCurrenciesRoutes(api, controller)
	s.gin = gin
	s.repoMock = repo
}

func (s *CurrenciesControllerTestSuite) TestGetLastBtcUsdtCourse() {
	course := db.BtcUsdtCourseSlice{Value: 100, Timestamp: time.Now()}
	s.repoMock.On("GetLastBtcUsdt").Return(&course, nil)
	req, err := http.NewRequest("GET", "/api/btcusdt", nil)
	s.NoError(err)
	resp := httptest.NewRecorder()
	s.gin.ServeHTTP(resp, req)
	s.Equal(200, resp.Code)
	var respBody rest.BtcUsdtCourseSlice
	err = json.Unmarshal(resp.Body.Bytes(), &respBody)
	s.NoError(err)
	s.Equal(course.Value, respBody.Value)
	s.Equal(uint64(course.Timestamp.Unix()), respBody.Timestamp)
}

func (s *CurrenciesControllerTestSuite) TestGetBtcUsdtCourseHistory() {
	time := nowWithoutNanos()
	course := db.BtcUsdtCourseSlice{Value: 100, Timestamp: time}
	s.repoMock.
		On("GetBtcUsdtHistory", uint(0), uint(0), time).
		Return(int64(1), []db.BtcUsdtCourseSlice{course}, nil)
	pagination := rest.BtcUsdtPaginationRequest{
		Limit: 0, Offset: 0,
		FromTime: time.Unix(),
	}
	reqBody, err := json.Marshal(pagination)
	s.NoError(err)
	req, err := http.NewRequest("POST", "/api/btcusdt", bytes.NewBuffer(reqBody))
	s.NoError(err)
	resp := httptest.NewRecorder()
	s.gin.ServeHTTP(resp, req)
	s.Equal(200, resp.Code)
	var respBody rest.BtcUsdtHistoryPage
	err = json.Unmarshal(resp.Body.Bytes(), &respBody)
	s.NoError(err)
	s.Equal(1, respBody.Total)
	s.Equal(1, len(respBody.History))
	s.Equal(course.Value, respBody.History[0].Value)
	s.Equal(uint64(course.Timestamp.Unix()), respBody.History[0].Timestamp)
}

func (s *CurrenciesControllerTestSuite) TestGetLastRubToFiatCourse() {
	course := db.AnyToFiatCourseSlice{
		BaseCurrencyCode: "RUB",
		FiatCode:         "USD",
		Value:            100,
		GetDate:          time.Now(),
	}
	s.repoMock.
		On("GetLastAnyToFiat", "RUB").
		Return([]db.AnyToFiatCourseSlice{course}, nil)
	req, err := http.NewRequest("GET", "/api/currencies", nil)
	s.NoError(err)
	resp := httptest.NewRecorder()
	s.gin.ServeHTTP(resp, req)
	s.Equal(200, resp.Code)
	var respBody map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &respBody)
	s.NoError(err)
	s.Equal(course.BaseCurrencyCode, respBody["base_currency"].(string))
	s.Equal(course.Value, respBody[course.FiatCode].(float64))
}

func (s *CurrenciesControllerTestSuite) TestGetRubToFiatCourseHistory() {
	date := rest.DateOnly(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)) //fixme
	course := db.AnyToFiatCourseSlice{
		BaseCurrencyCode: "RUB",
		FiatCode:         "USD",
		Value:            100,
		GetDate:          time.Time(date),
	}
	s.repoMock.
		On("GetAnyToFiatHistory", "RUB", uint(0), uint(0), time.Time(date)).
		Return(int64(1), []db.AnyToFiatCourseSlice{course}, nil)
	pagination := rest.AnyToFiatPaginationRequest{
		Limit: 0, Offset: 0,
		FromDate: date,
	}
	reqBody, err := json.Marshal(pagination)
	s.NoError(err)
	req, err := http.NewRequest("POST", "/api/currencies", bytes.NewBuffer(reqBody))
	s.NoError(err)
	resp := httptest.NewRecorder()
	s.gin.ServeHTTP(resp, req)
	s.Equal(200, resp.Code)
}

func (s *CurrenciesControllerTestSuite) TestGetLastBtcToFiatCourse() {
	course := db.AnyToFiatCourseSlice{
		BaseCurrencyCode: "BTC",
		FiatCode:         "USD",
		Value:            100,
		GetDate:          time.Now(),
	}
	s.repoMock.
		On("GetLastAnyToFiat", "BTC").
		Return([]db.AnyToFiatCourseSlice{course}, nil)
	req, err := http.NewRequest("GET", "/api/latest", nil)
	s.NoError(err)
	resp := httptest.NewRecorder()
	s.gin.ServeHTTP(resp, req)
	s.Equal(200, resp.Code)
	var respBody map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &respBody)
	s.NoError(err)
	s.Equal(course.BaseCurrencyCode, respBody["base_currency"].(string))
	s.Equal(course.Value, respBody[course.FiatCode].(float64))	
}

func nowWithoutNanos() time.Time {
	t := time.Now()
	return time.Date(t.Day(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
}

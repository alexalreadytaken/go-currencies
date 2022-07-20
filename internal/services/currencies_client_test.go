package services

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
)

type CurrenciesClientTestSuite struct {
	suite.Suite
	client *CurrenciesClient
}

func TestCurrenciesClient(t *testing.T) {
	suite.Run(t, &CurrenciesClientTestSuite{})
}

const (
	coinsUrl      = "http://localhost:2121"
	coinsResponse = `{
		"rate":100.12
	}`
	fiatUrl      = "http://localhost:3131"
	fiatResponse = `{
		"base":"RUB",
		"rates":{
			"USD":1203.09
		}
	}`
)

func (s *CurrenciesClientTestSuite) SetupSuite() {
	httpmock.Activate()
	resty := resty.New()
	httpmock.ActivateNonDefault(resty.GetClient())

	httpmock.RegisterResponder("GET", fiatUrl+"/?"+rubToFiatQuery,
		httpmock.NewStringResponder(200, fiatResponse))

	httpmock.RegisterResponder("POST", coinsUrl+"/",
		httpmock.NewStringResponder(200, coinsResponse))

	s.client = &CurrenciesClient{
		resty:                  resty,
		coinsApiKey:            "",
		coinsApiUrl:            coinsUrl,
		coinsApiBtcUsdEndpoint: "",
		fiatApiUrl:             fiatUrl,
		fiatApiPriceEndpoint:   "",
	}
}

func (s *CurrenciesClientTestSuite) TestPositiveResponsesToClient() {
	course, err := s.client.GetBtcUsdCourse()
	s.NoError(err)
	s.Equal(100.12, course)

	courses, err := s.client.GetRubToFiatCourse()
	s.NoError(err)
	s.Equal("RUB", courses.Base)
	s.Equal(1203.09, courses.Rates["USD"])
}

func (s *CurrenciesClientTestSuite) TearDownSuite() {
	httpmock.DeactivateAndReset()
}

package services

import (
	"encoding/json"
	"fmt"

	"github.com/alexalreadytaken/go-currencies/internal/models/rest"
	"github.com/alexalreadytaken/go-currencies/internal/utils"
	"github.com/go-resty/resty/v2"
)

type CurrenciesClient struct {
	resty                  *resty.Client
	coinsApiKey            string
	coinsApiUrl            string
	coinsApiBtcUsdEndpoint string
	fiatApiUrl             string
	fiatApiPriceEndpoint   string
}

func NewCurrenciesClient(cnf *utils.AppConfig) *CurrenciesClient {
	return &CurrenciesClient{
		resty:                  resty.New(),
		coinsApiKey:            cnf.CoinsClientApiKey,
		coinsApiUrl:            cnf.CoinsClientApiUrl,
		coinsApiBtcUsdEndpoint: cnf.CoinsClientApiBtcUsdEndpoint,
		fiatApiUrl:             cnf.FiatClientApiUrl,
		fiatApiPriceEndpoint:   cnf.FiatClientApiPriceEndpoint,
	}
}

const rubToFiatQuery = "base=RUB&symbols=CNY,USD,EUR,JPY,GBP,KPW,INR,CAD,HKD,TWD,AUD,BRL,CHF,THB,MXN,SAR,SGD"

var btcUsdtRequest = rest.CoinPriceRequest{
	Currecny: "USDT",
	Code:     "BTC",
	Meta:     false,
}

func (client *CurrenciesClient) GetBtcUsdCourse() (float64, error) {
	resp, err := client.resty.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", client.coinsApiKey).
		SetBody(btcUsdtRequest).
		Post(fmt.Sprintf("%s/%s", client.coinsApiUrl, client.coinsApiBtcUsdEndpoint))
	if err != nil {
		return 0, err
	}
	if resp.StatusCode() != 200 {
		return 0, fmt.Errorf("api server return invalid code, response= %+v", resp.RawResponse)
	}
	var respBody rest.CoinPriceResponse
	err = json.Unmarshal(resp.Body(), &respBody)
	if err != nil {
		return 0, fmt.Errorf("can't parse reponse json = %s", err.Error())
	}
	return respBody.Rate, nil
}

func (client *CurrenciesClient) GetRubToFiatCourse() (*rest.FiatPricesResposnse, error) {
	resp, err := client.resty.R().
		SetHeader("Content-Type", "application/json").
		Get(fmt.Sprintf("%s/%s?%s",
			client.fiatApiUrl, client.fiatApiPriceEndpoint, rubToFiatQuery))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("api server return invalid code, response= %+v", resp.RawResponse)
	}
	var respBody rest.FiatPricesResposnse
	err = json.Unmarshal(resp.Body(), &respBody)
	if err != nil {
		return nil, fmt.Errorf("can't parse reponse json = %s", err.Error())
	}
	return &respBody, nil
}

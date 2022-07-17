package utils

import (
	"log"
	"os"
)

//todo to config lib
type AppConfig struct {
	DbUser                       string
	DbPass                       string
	DbHost                       string
	DbPort                       string
	DbName                       string
	CoinsClientApiKey            string
	CoinsClientApiUrl            string
	CoinsClientApiBtcUsdEndpoint string
	FiatClientApiUrl             string
	FiatClientApiPriceEndpoint   string
}

func LoadConfigFromEnv() *AppConfig {
	return &AppConfig{
		DbUser:                       loadEnvByKey("DB_USER"),
		DbPass:                       loadEnvByKey("DB_PASS"),
		DbHost:                       loadEnvByKey("DB_HOST"),
		DbPort:                       loadEnvByKey("DB_PORT"),
		DbName:                       loadEnvByKey("DB_NAME"),
		CoinsClientApiKey:            loadEnvByKey("COINS_API_KEY"),
		CoinsClientApiUrl:            loadEnvByKey("COINS_API_URL"),
		CoinsClientApiBtcUsdEndpoint: loadEnvByKey("COINS_API_COINPRICE_ENDPOINT"),
		FiatClientApiUrl:             loadEnvByKey("FIAT_API_URL"),
		FiatClientApiPriceEndpoint:   loadEnvByKey("FIAT_PRICE_ENDPOINT"),
	}
}

func loadEnvByKey(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("can't load env value by key = %s", key)
	}
	return val
}

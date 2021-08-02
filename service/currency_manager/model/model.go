package model

import "strings"

type CurrencyType string

const (
	USD CurrencyType = "USD"
	EUR CurrencyType = "EUR"
	AUD CurrencyType = "AUD"
	CAD CurrencyType = "CAD"
	GBP CurrencyType = "GBP"
	PLN CurrencyType = "PLN"
	UAH CurrencyType = "UAH"
)

var AllCurrenciesTypes = map[CurrencyType]struct{}{
	USD: {},
	EUR: {},
	AUD: {},
	CAD: {},
	GBP: {},
	PLN: {},
	UAH: {},
}

func AllCurrenciesTypesString() string {
	var currencies string
	for cur, _ := range AllCurrenciesTypes {
		currencies += string(cur) + ","
	}
	currencies = strings.TrimSuffix(currencies, ",")
	return currencies
}

type CurrencyEntry struct {
	Type  CurrencyType
	Ratio float64
}

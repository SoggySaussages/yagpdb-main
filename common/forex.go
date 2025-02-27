package common

import (
	"encoding/json"
	"errors"
	"fmt"

	"net/http"
)

const CurrenciesAPIURL = "https://api.frankfurter.app/currencies"

var ErrUnknownForexCurrencyError = errors.New("unknown currency")
var ErrFailedConversion = errors.New("failed to convert")

type ErrUnknownForexCurrency struct {
	Currency string
}

func (e *ErrUnknownForexCurrency) Error() string {
	return "Unknown/Unsupported currency: " + e.Currency
}

func (e *ErrUnknownForexCurrency) Is(target error) bool {
	return errors.Is(target, ErrUnknownForexCurrencyError)
}

func ForexConvert(amount float64, from, to string) (*Currencies, *ExchangeRate, error) {
	var currenciesResult Currencies
	var exchangeRateResult ExchangeRate

	err := forexRequestAPI(CurrenciesAPIURL, &currenciesResult)
	if err != nil {
		return nil, nil, err
	}

	// Check if the currencies exist in the map
	_, fromExist := currenciesResult[from]
	_, toExist := currenciesResult[to]

	// If the currency isn't supported by API.
	if !fromExist {
		return nil, nil, &ErrUnknownForexCurrency{from}
	} else if !toExist {
		return nil, nil, &ErrUnknownForexCurrency{to}
	}

	err = forexRequestAPI(fmt.Sprintf("https://api.frankfurter.app/latest?amount=%.3f&from=%s&to=%s", amount, from, to), &exchangeRateResult)
	if err != nil {
		return nil, nil, err
	}

	return &currenciesResult, &exchangeRateResult, nil
}

func forexRequestAPI(query string, result interface{}) error {
	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "SGPDB.xyz (https://github.com/botlabs-gg/sgpdb)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrFailedConversion
	}

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return err
	}
	return nil
}

type ExchangeRate struct {
	Amount float64
	Base   string
	Date   string
	Rates  map[string]float64
}
type Currencies map[string]string

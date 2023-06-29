package exchangeRate

import (
	"fmt"
	"io"
	"log"
	"time"
	"encoding/json"
	"net/http"
)

type Currency map[string]float64

type Rates map[string]Currency

type ExchangeRates struct {
	Rates
	Interval time.Duration
}

func New(interval time.Duration) *ExchangeRates {
	return &ExchangeRates{ 
		Rates: Rates{}, 
		Interval: interval,
	}
}

func (e *ExchangeRates) FetchExchangeRates() {
	tokens := []string{
		"bitcoin",
		"ethereum",
		"litecoin",
	}
	currency := []string{
		"usd",
		"eur",
		"gbp",
	}
	var tokenString string
	var currencyString string

	for _, item := range tokens {
		tokenString += item + "," 
	}

	for _, item := range currency {
		currencyString += item + ","
	}

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=%s", tokenString, currencyString)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return
	}
	err = json.Unmarshal(body, &e.Rates)
	if err != nil {
		log.Fatalln(err)
		return
	}
}
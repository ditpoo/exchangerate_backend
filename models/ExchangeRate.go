package models

import (
	"time"
)

type Price struct {
	Cryptocurrency string    `json:"cryptocurrency" bson:"cryptocurrency"`
	FiatCurrency   string    `json:"fiatCurrency" bson:"fiatCurrency"`
	Price          float64   `json:"price" bson:"price"`
	Timestamp      time.Time `json:"timestamp" bson:"timestamp"`
}

type ExchangeRate struct {
	Rates map[string][]Price    `json:"rates" bson:"rates"`
	Timestamp   time.Time       `json:"timestamp" bson:"timestamp"`
	Type string                 `josn:"type" bosn:"type"`
}
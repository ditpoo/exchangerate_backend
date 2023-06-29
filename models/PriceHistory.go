package models

import (
	"time"
)

type PriceHistory struct {
	Cryptocurrency string    `json:"cryptocurrency" bson:"cryptocurrency"`
	FiatCurrency   string    `json:"fiatCurrency" bson:"fiatCurrency"`
	Price          float64   `json:"price" bson:"price"`
	Timestamp      time.Time `json:"timestamp" bson:"timestamp"`
}
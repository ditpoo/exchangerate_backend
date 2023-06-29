package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"firebond/models"
	"firebond/pkg/exchangeRate"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) InsertCurrentPrice(currentExchangeRates exchangeRate.Rates) {
	cryptoMap := make(map[string][]models.Price)
	for cryptoToken, rates := range currentExchangeRates {
		var fiatList []models.Price
		for fiatCurrency, price := range rates {
			p := models.Price{
				FiatCurrency: fiatCurrency,
				Cryptocurrency: cryptoToken,
				Price: price,
				Timestamp: time.Now(),
			}
			fiatList = append(fiatList, p)
		}
		cryptoMap[cryptoToken] = fiatList
	}
	filter := bson.M{
		"type": "CURRENT_EXCHANGE_RATE",
	}
	updateOptions := options.Update().SetUpsert(true)
	updateResult, err := r.db.Collection("current_rates").UpdateOne(
		context.TODO(), 
		filter, bson.M{"$set": models.ExchangeRate{
			Rates: cryptoMap,
			Timestamp: time.Now(),
			Type: "CURRENT_EXCHANGE_RATE",
		}}, 
		updateOptions)
	if err != nil {
		log.Fatal(err)
	}

	if updateResult.UpsertedID != nil {
		log.Println("Document inserted:", updateResult.UpsertedID)
	} else {
		log.Println("Document updated:", filter)
	}
}

func (r *Repository) InsertHistoricalPrice(currentExchangeRates exchangeRate.Rates) {
	var fiatList []interface{}
	for cryptoToken, rates := range currentExchangeRates {
		for fiatCurrency, price := range rates {
			p := models.Price{
				FiatCurrency: fiatCurrency,
				Cryptocurrency: cryptoToken,
				Price: price,
				Timestamp: time.Now(),
			}
			fiatList = append(fiatList, p)
		}
	}
	insertOptions := options.InsertMany().SetOrdered(false)
	insertResult, err := r.db.Collection("price_history").InsertMany(context.TODO(), fiatList, insertOptions)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Inserted document IDs: %+v", insertResult.InsertedIDs)
	log.Printf("\n Successfully Inserted Historical Prices %+v", currentExchangeRates)
}

func (r *Repository) GetCurrentRates(cryptoToken string, fiatCurrency string) (models.PriceHistory, error) {
	filter := bson.M{
		"cryptocurrency": cryptoToken,
		"fiatCurrency":   fiatCurrency,
	}
	opts := options.FindOne().SetSort(bson.M{"timestamp": -1})
	var result models.PriceHistory
	err := r.db.Collection("price_history").FindOne(context.Background(), filter, opts).Decode(&result)	
	if err != nil {
		return result, err
	}
	return result, nil
}

func (r *Repository) GetAllCurrentRates(cryptoToken string) ([]models.Price, error) {
	filter := bson.M{
		"type": "CURRENT_EXCHANGE_RATE",
	}
	opts := options.FindOne().SetSort(bson.M{"timestamp": -1})
	var result models.ExchangeRate
	err := r.db.Collection("current_rates").FindOne(context.Background(), filter, opts).Decode(&result)	
	if err != nil {
		return []models.Price{}, err
	}
	for key, val := range result.Rates {
		if key == cryptoToken {
			return val, nil
		}
	}
	return []models.Price{}, errors.New("not found")
}

func (r *Repository) GetAllExchangeRates() (models.ExchangeRate, error) {
	filter := bson.M{
		"type": "CURRENT_EXCHANGE_RATE",
	}
	opts := options.FindOne().SetSort(bson.M{"timestamp": -1})
	var exchangeRates models.ExchangeRate
	err := r.db.Collection("current_rates").FindOne(context.Background(), filter, opts).Decode(&exchangeRates)	
	if err != nil {
		return exchangeRates, err
	}
	return exchangeRates, nil
}

func (r *Repository) GetHistoricalRates(cryptoToken string, fiatCurrency string) ([]models.PriceHistory, error) {
	twentyFourHoursAgo := time.Now().Add(-24 * time.Hour)
	filter := bson.M{
		"cryptocurrency": cryptoToken,
		"fiatCurrency":   fiatCurrency,
		"timestamp": bson.M{
			"$gte": twentyFourHoursAgo,
		},
	}
	opts := options.Find().SetSort(bson.M{"timestamp": 1})
	cur, err := r.db.Collection("price_history").Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	var results []models.PriceHistory
	for cur.Next(context.Background()) {
		var result models.PriceHistory
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}
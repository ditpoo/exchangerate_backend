package main

import (
	"net/http"
	"os"
	"time"

	"context"
	"log"
	"sync"

	"firebond/handlers"
	"firebond/pkg/exchangeRate"
	"firebond/repository"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	e := exchangeRate.New(5 * time.Second)
	e.FetchExchangeRates()
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URL"))
	client, err := mongo.Connect(context.TODO(), clientOptions);
	if err != nil {
		log.Fatal(err)
	}
	repository := repository.NewRepository(client.Database("app"))
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for {
			e.FetchExchangeRates()
			log.Printf("\n %+v \n", e.Rates)
			go repository.InsertCurrentPrice(e.Rates)
			go repository.InsertHistoricalPrice(e.Rates)
			time.Sleep(e.Interval)
		}
	}()

	exchangeHandler := handlers.ExchangeRateHandler{
		Repository: repository,
	}

	router := mux.NewRouter()
	router.HandleFunc("/health", handlers.HealthHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/rates/{cryptocurrency}/{fiat}", exchangeHandler.GetCurrentRate).Methods("GET", "OPTIONS")
	router.HandleFunc("/rates/{cryptocurrency}", exchangeHandler.GetAllRates).Methods("GET", "OPTIONS")
	router.HandleFunc("/rates", exchangeHandler.GetAllExchangeRates).Methods("GET", "OPTIONS")
	router.HandleFunc("/rates/history/{cryptocurrency}/{fiat}", exchangeHandler.GetHistoricalRates).Methods("GET", "OPTIONS")
	router.HandleFunc("/balance/{address}", exchangeHandler.GetBalance).Methods("GET", "OPTIONS")
	
	http.Handle("/", router)

	http.ListenAndServe(":5000", router)
	wg.Wait()
}
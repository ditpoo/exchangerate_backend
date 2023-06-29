package handlers

import (
	// "log"
	"log"
	"net/http"

	"firebond/common"
	"firebond/repository"
	"firebond/pkg/web3"

	"github.com/gorilla/mux"
)

type ExchangeRateHandler struct {
	Repository *repository.Repository
}

func (r *ExchangeRateHandler) GetCurrentRate(resp http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	cryptocurrency := params["cryptocurrency"]
	fiatCurrency := params["fiat"]
	log.Println(cryptocurrency, common.IsValidCryptoToken(cryptocurrency))
	if !common.IsValidCryptoToken(cryptocurrency) {
		common.BadRequestResponse("Invalid cryptocurrency param", resp)
		return
	}
	if !common.IsValidFiat(fiatCurrency) {
		common.BadRequestResponse("Invalid fiat param", resp)
		return
	}
	result, err := r.Repository.GetCurrentRates(cryptocurrency, fiatCurrency)
	if err != nil {
		log.Println(err)
		common.InternalErrorResponse("Failed To fetch Exchange Rate", resp)
		return
	}
	common.SuccessFullResponse(result, resp)
}

func (r *ExchangeRateHandler) GetAllRates(resp http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	cryptocurrency := params["cryptocurrency"]
	if !common.IsValidCryptoToken(cryptocurrency) {
		common.BadRequestResponse("Invalid cryptocurrency param", resp)
		return
	}
	result, err := r.Repository.GetAllCurrentRates(cryptocurrency)
	if err != nil {
		log.Println(err)
		common.InternalErrorResponse("Failed To fetch Exchange Rate", resp)
		return
	}
	common.SuccessFullResponse(result, resp)
}

func (r *ExchangeRateHandler) GetAllExchangeRates(resp http.ResponseWriter, req *http.Request) {
	result, err := r.Repository.GetAllExchangeRates()
	if err != nil {
		log.Println(err)
		common.InternalErrorResponse("Failed To fetch Exchange Rate", resp)
		return
	}
	common.SuccessFullResponse(result, resp)
}

func (r *ExchangeRateHandler) GetHistoricalRates(resp http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	cryptocurrency := params["cryptocurrency"]
	fiatCurrency := params["fiat"]
	if !common.IsValidCryptoToken(cryptocurrency) {
		common.BadRequestResponse("Invalid cryptocurrency param", resp)
		return
	}
	if !common.IsValidFiat(fiatCurrency) {
		common.BadRequestResponse("Invalid fiat param", resp)
		return
	}
	result, err := r.Repository.GetHistoricalRates(cryptocurrency, fiatCurrency)
	if err != nil {
		log.Println(err)
		common.InternalErrorResponse("Failed To fetch historical Exchange Rates", resp)
		return
	}
	common.SuccessFullResponse(result, resp)
}

func (r *ExchangeRateHandler) GetBalance(resp http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	address := params["address"]
	if !common.IsValidEthereumAddress(address) {
		common.BadRequestResponse("Invalid ETH Address", resp)
		return
	}
	w := web3.New(address)
	w.UpdateBalance()
	common.SuccessFullResponse(w, resp)
}
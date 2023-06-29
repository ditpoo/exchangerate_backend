package common

import (
	"encoding/json"
	"net/http"
	"regexp"
)

type ErrorResponse struct {
	Error string       `json:"error"`
	StatusCode  int    `json:"status"`
	Mesage string      `json:"message"`
}

type Response struct {
	Message string    `json:"message"`
	StatusCode int    `json:"status"`
	Data interface{}  `json:"data"`
}

func GenerateResponse(err string, msg string, status int, w http.ResponseWriter) {
	var response = ErrorResponse{
		Error:  err,
		Mesage: msg,
		StatusCode:  status,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

func BadRequestResponse(msg string, w http.ResponseWriter) {
	GenerateResponse("Bad Request", msg, http.StatusBadRequest, w)
}

func InternalErrorResponse(msg string, w http.ResponseWriter) {
	GenerateResponse("Internal Server Error", msg, http.StatusInternalServerError, w)
}

func SuccessFullResponse(data interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Message: "Successfull",
		StatusCode: http.StatusOK,
		Data: data,
	})
}

func IsValidCryptoToken(token string) bool {
	tokenMap := map[string]string{
		"bitcoin": "bitcoin",
		"ethereum": "ethereum",
		"litecoin": "litecoin",
	}
	_, ok := tokenMap[token]
	return ok
}

func IsValidFiat(fiat string) bool {
	fiatMap := map[string]string{
		"usd": "usd",
		"eur": "eur",
		"gbp": "gbp",
	}
	_, ok := fiatMap[fiat]
	return ok
}

func IsValidEthereumAddress(address string) bool {
	// Check if hex and length is 40
	match, _ := regexp.MatchString("^0x[0-9a-fA-F]{40}$", address)
	// and not zeros
	return match && address != "0x0000000000000000000000000000000000000000"
}

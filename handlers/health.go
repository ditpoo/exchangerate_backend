package handlers

import (
	"net/http"
)

func HealthHandler(resp http.ResponseWriter, req *http.Request) {	
	resp.WriteHeader(http.StatusOK)
}
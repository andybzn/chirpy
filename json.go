package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func returnError(w http.ResponseWriter, statusCode int, message string, err error) {
	if err != nil {
		log.Printf("Error: %v", err)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}

	res, err := json.Marshal(errorResponse{Error: message})
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(res)
}

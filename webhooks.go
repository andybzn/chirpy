package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

func (cfg *apiConfig) handlerUserUpgrade(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId string `json:"user_id"`
		} `json:"data"`
	}

	// get the request parameters
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		returnError(w, http.StatusInternalServerError, "Error unmarshalling JSON", err)
		return
	}

	// ignore non-`user.upgraded` events
	if params.Event != "user.upgraded" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte{})
		return
	}

	// parse userId
	userId, err := uuid.Parse(params.Data.UserId)
	if err != nil {
		returnError(w, http.StatusInternalServerError, "Error parsing user_id", err)
		return
	}

	// upgrade user
	if err := cfg.db.UpgradeToChirpyRedById(r.Context(), userId); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte{})
}

package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		returnError(w, http.StatusInternalServerError, "Error unmarshalling JSON", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		returnError(w, http.StatusInternalServerError, "Error creating user", err)
		return
	}

	data, err := json.Marshal(User{user.ID, user.CreatedAt, user.UpdatedAt, user.Email})
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		returnError(w, http.StatusInternalServerError, "Error marshalling JSON", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

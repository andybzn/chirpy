package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/andybzn/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	type returnValue struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		returnError(w, http.StatusInternalServerError, "Error unmarshalling JSON", err)
		return
	}

	// validate the chirp
	chirpBody, err := validateChirp(params.Body)
	if err != nil {
		returnError(w, http.StatusBadRequest, "Bad request", err)
	}

	// save the chirp
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   chirpBody,
		UserID: params.UserId,
	})
	if err != nil {
		log.Printf("Error creating chirp: %v", err)
		returnError(w, http.StatusInternalServerError, "Error creating chirp", err)
		return
	}

	data, err := json.Marshal(returnValue{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		returnError(w, http.StatusInternalServerError, "Error marshalling JSON", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func validateChirp(chirp string) (string, error) {
	const maxLength = 140
	if len(chirp) > maxLength {
		return "", errors.New("Chirp is too long")
	}

	profanity := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	return replaceProfanity(profanity, chirp), nil
}

func replaceProfanity(profanity map[string]struct{}, text string) string {
	splits := strings.Split(text, " ")
	for i, word := range splits {
		if _, ok := profanity[strings.ToLower(word)]; ok {
			splits[i] = "****"
		}
	}

	return strings.Join(splits, " ")
}

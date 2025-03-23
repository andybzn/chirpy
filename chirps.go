package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/andybzn/chirpy/internal/auth"
	"github.com/andybzn/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnValue struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    uuid.UUID `json:"user_id"`
	}

	// validate the user
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		returnError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), nil)
		return
	}
	userId, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		returnError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), nil)
		return
	}

	// validate the request
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
		UserID: userId,
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

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	cheeps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		log.Printf("Error getting chirps from database: %v", err)
		returnError(w, http.StatusInternalServerError, "Error getting chirps", err)
		return
	}

	chirps := []Chirp{}
	for _, chirp := range cheeps {
		chirps = append(chirps, Chirp{
			Id:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		})
	}

	data, err := json.Marshal(chirps)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		returnError(w, http.StatusInternalServerError, "Error marshalling JSON", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (cfg *apiConfig) handlerGetChirpsById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpId")
	if id == "" {
		returnError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), nil)
		return
	}
	parsed_id, err := uuid.Parse(id)
	if err != nil {
		returnError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), nil)
		return
	}

	data, err := cfg.db.GetChirpsById(r.Context(), parsed_id)
	if err != nil {
		returnError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), nil)
		return
	}

	chirp, err := json.Marshal(Chirp{
		Id:        data.ID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		Body:      data.Body,
		UserId:    data.UserID,
	})
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		returnError(w, http.StatusInternalServerError, "Error marshalling JSON", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(chirp)
}

func (cfg *apiConfig) handlerDeleteChirpsById(w http.ResponseWriter, r *http.Request) {
	// PATH VALIDATION
	id := r.PathValue("chirpId")
	if id == "" {
		returnError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), nil)
		return
	}
	// USER VALIDATION
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		returnError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), nil)
		return
	}
	userId, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		returnError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), nil)
		return
	}
	// CHIRP VALIDATION
	chirpId, err := uuid.Parse(id)
	if err != nil {
		returnError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), nil)
		return
	}
	chirp, err := cfg.db.GetChirpsById(r.Context(), chirpId)
	if err != nil {
		returnError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), nil)
		return
	}

	// Verify that the calling user is the author of the chirp
	if chirp.UserID != userId {
		returnError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden), nil)
		return
	}

	// Delete the chirp
	if err := cfg.db.DeleteChirpsById(r.Context(), database.DeleteChirpsByIdParams{ID: chirpId, UserID: userId}); err != nil {
		returnError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(http.StatusText(http.StatusNoContent)))
}

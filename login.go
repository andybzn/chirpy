package main

import (
	"encoding/json"
	"github.com/andybzn/chirpy/internal/database"
	"log"
	"net/http"
	"time"

	"github.com/andybzn/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type responseData struct {
		User
		AccessToken  string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		returnError(w, http.StatusInternalServerError, "Error unmarshalling JSON", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		returnError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	if err := auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil {
		returnError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	// JSON WEB TOKENS
	// get the token
	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		returnError(w, http.StatusInternalServerError, "Failed to generate token", err)
		return
	}
	// create a refresh token
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		returnError(w, http.StatusInternalServerError, "Failed to generate token", err)
		return
	}
	// store the refresh token
	_, err = cfg.db.StoreRefreshToken(r.Context(), database.StoreRefreshTokenParams{
		Token:     refreshToken,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().AddDate(0, 0, 60),
	})
	if err != nil {
		returnError(w, http.StatusInternalServerError, "Failed to store token", err)
		return
	}

	// return the user object
	data, err := json.Marshal(responseData{User{user.ID, user.CreatedAt, user.UpdatedAt, user.Email}, accessToken, refreshToken})
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		returnError(w, http.StatusInternalServerError, "Error marshalling JSON", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

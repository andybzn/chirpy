package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/andybzn/chirpy/internal/auth"
	"github.com/andybzn/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		returnError(w, http.StatusInternalServerError, "Error unmarshalling JSON", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		returnError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		returnError(w, http.StatusInternalServerError, "Error creating user", err)
		return
	}

	data, err := json.Marshal(User{user.ID, user.CreatedAt, user.UpdatedAt, user.Email, user.IsChirpyRed})
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		returnError(w, http.StatusInternalServerError, "Error marshalling JSON", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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
	if params.Email == "" || params.Password == "" {
		returnError(w, http.StatusBadRequest, "`email` and `password` fields must be provided", err)
		return
	}

	// update the user
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		returnError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}
	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{ID: userId, Email: params.Email, HashedPassword: hashedPassword})
	if err != nil {
		returnError(w, http.StatusInternalServerError, "Error updating user", err)
		return
	}

	// return a response
	data, err := json.Marshal(User{user.ID, user.CreatedAt, user.UpdatedAt, user.Email, user.IsChirpyRed})
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		returnError(w, http.StatusInternalServerError, "Error marshalling JSON", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

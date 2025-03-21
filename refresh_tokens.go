package main

import (
	"encoding/json"
	"github.com/andybzn/chirpy/internal/auth"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type responseData struct {
		AccessToken string `json:"token"`
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		returnError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), nil)
		return
	}
	tokenDetails, err := cfg.db.GetUserByRefreshToken(r.Context(), bearerToken)
	if err != nil {
		returnError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), nil)
		return
	}
	if tokenDetails.ExpiresAt.Before(time.Now().UTC()) || tokenDetails.RevokedAt.Valid == true {
		returnError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), nil)
		return
	}
	accessToken, err := auth.MakeJWT(tokenDetails.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		returnError(w, http.StatusInternalServerError, "Failed to generate token", err)
		return
	}
	data, err := json.Marshal(responseData{accessToken})
	if err != nil {
		returnError(w, http.StatusInternalServerError, "Failed to generate token", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		returnError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), nil)
		return
	}
	tokenDetails, err := cfg.db.GetUserByRefreshToken(r.Context(), bearerToken)
	if err != nil {
		returnError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), nil)
		return
	}
	if tokenDetails.ExpiresAt.Before(time.Now().UTC()) || tokenDetails.RevokedAt.Valid == true {
		returnError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), nil)
		return
	}

	err = cfg.db.RevokeToken(r.Context(), bearerToken)
	if err != nil {
		returnError(w, http.StatusInternalServerError, "Error revoking token", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(http.StatusText(http.StatusNoContent)))
}

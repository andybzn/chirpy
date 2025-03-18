package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnValue struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		returnError(w, http.StatusInternalServerError, "Error unmarshalling JSON", err)
		return
	}

	const maxLength = 140
	if len(params.Body) > maxLength {
		returnError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleaned := replaceProfanity([]string{"kerfuffle", "sharbert", "fornax"}, params.Body)

	data, err := json.Marshal(returnValue{CleanedBody: cleaned})
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		returnError(w, http.StatusInternalServerError, "Error marshalling JSON", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func replaceProfanity(profanity []string, text string) string {
	const replacement string = "****"
	splits := strings.Split(text, " ")
	for i, word := range splits {
		for _, p := range profanity {
			if strings.ToLower(word) == p {
				splits[i] = replacement
			}
		}
	}

	return strings.Join(splits, " ")
}

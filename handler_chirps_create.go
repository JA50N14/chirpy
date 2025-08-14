package main

import (
	"net/http"
	"encoding/json"
	"strings"
	"time"
	"errors"

	"github.com/JA50N14/chirpy/internal/auth"
	"github.com/JA50N14/chirpy/internal/database"
	"github.com/google/uuid"
)

const maxChirpLength = 140

var bannedWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert": {},
	"fornax": {},
}

type Chirp struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}


func (cfg *apiConfig) handlerChirpCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type response struct {
		Chirp
	}

	var params parameters

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not decode request", err)
		return
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "bearer token not in header", err)
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid bearer token", err)
		return
	}

	cleanedBody, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	params.Body = cleanedBody

	postChirpParams := database.PostChirpParams {
		Body: params.Body,
		UserID: userID,
	}


	chirp, err := cfg.db.PostChirp(r.Context(), postChirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error inserting chirp into database", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response {
		Chirp {
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		},
	})
}


func validateChirp(body string) (string, error) {
	if len(body) > maxChirpLength {
		return "", errors.New("chirp exceeds 140 characters.")
	}
	cleaned := cleanBody(body)
	return cleaned, nil
}


func cleanBody(body string) string {
	words := strings.Split(body, " ")
	for idx, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := bannedWords[loweredWord]; ok {
			words[idx] = "****"
		}
	}
	return strings.Join(words, " ")
}

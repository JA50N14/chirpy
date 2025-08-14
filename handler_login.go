package main

import (
	"net/http"
	"encoding/json"
	"time"

	"github.com/JA50N14/chirpy/internal/auth"
)


func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
		ExpiresInSeconds *int `json:"expires_in_seconds"`
	}

	type response struct {
		User
		Token string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	var params parameters

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not decode parameters", err)
		return
	}

	dbUser, err := cfg.db.GetUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	var expirationTime = time.Hour

	if params.ExpiresInSeconds != nil && *params.ExpiresInSeconds > 0 && *params.ExpiresInSeconds < 3600 {
		expirationTime = time.Second * time.Duration(*params.ExpiresInSeconds)
	}

	jwtToken, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, expirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error generating jwt token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID: dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email: dbUser.Email,
		},
		Token: jwtToken,
	})
}
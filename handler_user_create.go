package main

import (
	"net/http"
	"encoding/json"
	"time"

	"github.com/JA50N14/chirpy/internal/auth"
	"github.com/JA50N14/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
}


func (cfg *apiConfig) handlerUserCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
	}

	var params parameters

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding request", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error hashing password", err)
		return
	}
	params.Password = hashedPassword

	createUserParams := database.CreateUserParams{
		Email: params.Email,
		HashedPassword: params.Password,
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), createUserParams)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error creating user", err)
		return
	}

	user := User {
		ID: dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email: dbUser.Email,
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: user,
	})
}

package main

import (
	"net/http"

	"github.com/google/uuid"
)


func (cfg *apiConfig) handlerChirpRetrieve(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp ID", err)
		return
	}

	dbChirp, err := cfg.db.RetrieveChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "could not retrieve chirp.", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID: dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body: dbChirp.Body,
		UserID: dbChirp.UserID,
	})
}


func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.RetrieveChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not retrieve chirps", err)
		return
	}

	response := []Chirp{}

	for _, dbChirp := range dbChirps {
		chirp := Chirp {
			ID: dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body: dbChirp.Body,
			UserID: dbChirp.UserID,
		}
		response = append(response, chirp)
	}

	respondWithJSON(w, http.StatusOK, response)
}



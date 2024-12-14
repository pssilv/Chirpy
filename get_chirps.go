package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
  dbChirps, err := cfg.db.GetChirps(r.Context())
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve  chirps", err)
    return
  }

  chirps := []Chirp{}
 
  for _, chirp := range dbChirps {
    chirps = append(chirps, Chirp{
      ID: chirp.ID,
      CreatedAt: chirp.CreatedAt,
      UpdatedAt: chirp.UpdatedAt,
      Body: chirp.Body,
      UserID: chirp.UserID,
    })
  }
  

  respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handleChirpsGet(w http.ResponseWriter, r *http.Request) {

  id, err := uuid.Parse(r.PathValue("chirpID"))
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "Invalid chirp ID", err)
    return
  }
  
  dbChirp, err := cfg.db.GetChirp(r.Context(), id)
  if err != nil {
    respondWithError(w, 404, "Couldn't find chirp", err)
  }

  respondWithJSON(w, http.StatusOK, Chirp{
    ID: dbChirp.ID,
    CreatedAt: dbChirp.CreatedAt,
    UpdatedAt: dbChirp.UpdatedAt,
    Body: dbChirp.Body,
    UserID: dbChirp.UserID,
  })

}

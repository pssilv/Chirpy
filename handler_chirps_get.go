package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
  authorIDString := r.URL.Query().Get("author_id")
  sortQuery := r.URL.Query().Get("sort")

  authorID, err := uuid.Parse(authorIDString)
  if err != nil && sortQuery == "" {
    respondWithError(w, 400, "Invalid user ID", err)
    return
  }

  if sortQuery == "" {
    sortQuery = "asc"
  }

  dbChirps, err := cfg.db.GetChirps(r.Context())
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve  chirps", err)
    return
  }

  chirps := []Chirp{}

  if authorID != uuid.Nil {
    for _, chirp := range dbChirps {
      if authorID == chirp.UserID {
        chirps = append(chirps, Chirp{
          ID: chirp.ID,
          CreatedAt: chirp.CreatedAt,
          UpdatedAt: chirp.UpdatedAt,
          Body: chirp.Body,
          UserID: chirp.UserID,
        })
      }
    } 
  } else {
    for _, chirp := range dbChirps {
      chirps = append(chirps, Chirp{
        ID: chirp.ID,
        CreatedAt: chirp.CreatedAt,
        UpdatedAt: chirp.UpdatedAt,
        Body: chirp.Body,
        UserID: chirp.UserID,
      })
    }
  }

  if sortQuery == "desc" {
    sort.Slice(chirps, func(i, j int) bool {
      return chirps[j].CreatedAt.Before(chirps[i].CreatedAt)
    })
  }
  sort.Slice(chirps, func(i, j int) bool {
    return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
  })
  
  respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {

  chirpID, err := uuid.Parse(r.PathValue("chirpID"))
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "Invalid chirp ID", err)
    return
  }
  
  dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
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

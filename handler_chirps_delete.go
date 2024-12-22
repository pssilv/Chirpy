package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/pssilv/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
  token, err := auth.GetBearerToken(r.Header)
  if err != nil {
    respondWithError(w, 401, "Malformed token", err)
    return
  }

  userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
  if err != nil {
    respondWithError(w, 403, "Authentication failed", err)
    return
  }

  chirpID, err := uuid.Parse(r.PathValue("chirpID"))
  if err != nil {
    respondWithError(w, 400, "Invalid chirpID", err)
    return
  }

  chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
  if err != nil {
    respondWithError(w, 404, "Failed to get chirp", err)
    return
  }

  if chirp.UserID != userID {
    respondWithError(w, 403,"You can't delete this chirp", err)
    return
  }

  if err := cfg.db.DeleteChirp(r.Context(), chirp.ID); err != nil {
    respondWithError(w, 401, "Failed to delete chirp", err)
    return
  }

  w.WriteHeader(http.StatusNoContent)
}

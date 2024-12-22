package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/pssilv/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
  const upgradeEvent = "user.upgraded"

   type parameters struct {
    Event string `json:"event"`
    Data  struct {
      UserID uuid.UUID `json:"user_id"`
    }
  }

  ApiKey, err := auth.GetAPIKey(r.Header)
  if err != nil {
    respondWithError(w, http.StatusUnauthorized, "Failed to get API key", err)
  }
  if ApiKey != cfg.polkaKey {
    respondWithError(w, http.StatusUnauthorized, "You can't upgrade", err)
    return
  }

  decoder := json.NewDecoder(r.Body)
  params := parameters{}
  if err := decoder.Decode(&params); err != nil {
    respondWithError(w, http.StatusInternalServerError, err.Error(), err)
    return
  }

  if params.Event != upgradeEvent {
    w.WriteHeader(http.StatusNoContent)
    return
  }

  _, err = cfg.db.UpgradeToChirpyRed(r.Context(), params.Data.UserID)
  if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
      respondWithError(w, http.StatusNotFound, "Failed to find user", err)
      return
    }
    respondWithError(w, http.StatusInternalServerError, "Failed to upgrade user", err)
    return
  }

  w.WriteHeader(http.StatusNoContent)
}

package main

import (
	"net/http"
	"time"

	"github.com/pssilv/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
  type response struct {
    Token string `json:"token"`
  }

  tokenString, err := auth.GetBearerToken(r.Header)
  if err != nil {
    respondWithError(w, 401, "Invalid Authorization header", err)
    return
  }
 
  refreshToken, err := cfg.db.GetRefreshToken(r.Context(), tokenString)
  if err != nil {
    respondWithError(w, 401, "Invalid refresh token", err)
    return
  }
  if refreshToken.RevokedAt.Valid {
    respondWithError(w, 401, "Refresh Token got revoked", err)
    return
  }
  if time.Now().After(refreshToken.ExpiresAt) {
    respondWithError(w, 401, "Refresh token already expired", err)
    return
  }

  user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken.UserID)
  if err != nil {
    respondWithError(w, 401, "User with this refresh token doesn't exist", err)
    return
  }

  token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
  if err != nil {
    respondWithError(w, 401, "Couldn't generate token", err)
    return
  }

  respondWithJSON(w, 200, response{
    Token: token,
  })
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
  tokenString, err := auth.GetBearerToken(r.Header)
  if err != nil {
    respondWithError(w, 401, "Invalid Authorization header", err)
  }

  cfg.db.RevokeRefreshToken(r.Context(), tokenString)

  w.WriteHeader(http.StatusNoContent)
}

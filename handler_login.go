package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/pssilv/Chirpy/internal/auth"
	"github.com/pssilv/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
  type parameters struct {
    Email            string `json:"email"`
    Password         string `json:"password"`
  }

  type response struct {
    User
    Token        string `json:"token"`
    RefreshToken string `json:"refresh_token"`
  }

  decoder := json.NewDecoder(r.Body)
  params := parameters{}
  
  if err := decoder.Decode(&params); err != nil {
    respondWithError(w, http.StatusInternalServerError, "Couldn't decode", err)
    return
  }

  user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
  if err != nil {
    respondWithError(w, 401, "Incorrect email or password", err)
    return
  }

  if err := auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil {
    respondWithError(w, 401, "Incorrect email or password", err)
    return
  }

  acessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
  if err != nil {
    respondWithError(w, 401, "Couldn't create JWT", err)
    return
  }

  refreshToken, err := auth.MakeRefreshToken()
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
  }

  cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
    Token: refreshToken,
    UserID: user.ID,
    ExpiresAt: time.Now().UTC().Add(time.Hour * 1440),
  })

  respondWithJSON(w, http.StatusOK, response{
    User: User {
      ID: user.ID,
      CreatedAt: user.CreatedAt,
      UpdatedAt: user.UpdatedAt,
      Email: user.Email,
      IsChirpyRed: user.IsChirpyRed,
    },
    Token: acessToken,
    RefreshToken: refreshToken,
  })
}

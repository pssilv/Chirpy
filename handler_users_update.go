package main

import (
	"encoding/json"
	"net/http"

	"github.com/pssilv/Chirpy/internal/auth"
	"github.com/pssilv/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
  type parameters struct {
    Email    string `json:"email"`
    Password string `json:"password"`
  }

  type response struct {
    User
  }

  token, err := auth.GetBearerToken(r.Header)
  if err != nil {
    respondWithError(w, 401, "malformed token", err)
    return
  }

  userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
  if err != nil {
    respondWithError(w, 401, "Authentication failed", err)
    return
  }

  decoder := json.NewDecoder(r.Body)
  params := parameters{}
  decoder.Decode(&params)

  hashedPassword, err := auth.HashPassword(params.Password)
  if err != nil {
    respondWithError(w, 401, "Failed to hash password", err)
    return
  }

  updatedUser, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
    Email: params.Email,
    HashedPassword: hashedPassword,
    ID: userID,
  })
  if err != nil {
    respondWithError(w, 401, "Failed to update user", err)
  }

  respondWithJSON(w, http.StatusOK, response{
    User: User{
      ID: updatedUser.ID,
      CreatedAt: updatedUser.CreatedAt,
      UpdatedAt: updatedUser.UpdatedAt,
      Email: updatedUser.Email,
      IsChirpyRed: updatedUser.IsChirpyRed,
    },
  })  
}

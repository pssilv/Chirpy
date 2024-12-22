package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pssilv/Chirpy/internal/auth"
	"github.com/pssilv/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
  type paramaters struct {
    Email string `json:"email"`
    Password string `json:"password"`
  }

  type response struct {
    User
  }

  decoder := json.NewDecoder(r.Body)
  params := paramaters{}
  
  if err := decoder.Decode(&params); err != nil {
    respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
    return
  }

  hashedPassword, err := auth.HashPassword(params.Password)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "couldn't hash password", err)
  }

  userParams := database.CreateUserParams{
    Email: params.Email,
    HashedPassword: hashedPassword,
  }

  user, err := cfg.db.CreateUser(r.Context(), userParams)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "couldn't create user", err)
    return
  }

  respondWithJSON(w, 201, response {
    User: User {
      ID: user.ID,
      CreatedAt: user.CreatedAt,
      UpdatedAt: user.UpdatedAt,
      Email: user.Email,
      IsChirpyRed: user.IsChirpyRed,
    },
  })
}

type User struct {
  ID          uuid.UUID `json:"id"`
  CreatedAt   time.Time `json:"created_at"`
  UpdatedAt   time.Time `json:"updated_at"`
  Email       string    `json:"email"`
  IsChirpyRed bool      `json:"is_chirpy_red"`
}

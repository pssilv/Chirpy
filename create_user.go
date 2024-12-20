package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
  type paramaters struct {
    Email string `json:"email"`
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


  user, err := cfg.db.CreateUser(r.Context(), params.Email)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "couldn't create user", err)
    return
  }

  respondWithJSON(w, 201, response{
    User: User{
      ID: user.ID,
      CreatedAt: user.CreatedAt,
      UpdatedAt: user.UpdatedAt,
      Email: user.Email,
    },
  })
}

type User struct {
  ID        uuid.UUID `json:"id"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
  Email     string    `json:"email"`
}

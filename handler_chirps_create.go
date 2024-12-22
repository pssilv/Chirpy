package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pssilv/Chirpy/internal/auth"
	"github.com/pssilv/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {

  type paramaters struct {
    Body string `json:"body"`
  }

  authorization, err := auth.GetBearerToken(r.Header)

  if err != nil {
    respondWithError(w, 401, "Invalid header", err)
  }

  userID, err := auth.ValidateJWT(authorization, cfg.jwtSecret)
  if err != nil {
    respondWithError(w, 401, "Couldn't validate", err)
  }

  type response struct{
    Chirp
  }

  decoder := json.NewDecoder(r.Body)
  params := paramaters{}

  if err := decoder.Decode(&params); err != nil {
    respondWithError(w, http.StatusInternalServerError, err.Error(), err)
    return
  }

  cleaned, err := ValidateChirp(params.Body)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "Chirp is invalid", err)
  }

  chirpParams := database.CreateChirpParams{
    Body: cleaned,
    UserID: userID,
  }

  chirp, err := cfg.db.CreateChirp(r.Context(), chirpParams)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
    return
  }

  respondWithJSON(w, 201, response{
    Chirp: Chirp {
      ID: chirp.ID,
      CreatedAt: chirp.CreatedAt,
      UpdatedAt: chirp.UpdatedAt,
      Body: chirp.Body,
      UserID: chirp.UserID,
    },
  }) 
}

func ValidateChirp(body string) (string, error) {
  const maxChirpLength = 140

  if len(body) > maxChirpLength {
    return "", fmt.Errorf("Chirp size is longer than 140")
  }

  badWords := make(map[string]bool)
  badWords["kerfuffle"] = true 
  badWords["sharbert"] = true 
  badWords["fornax"] = true

  cleaned := getCleanedBody(body, badWords)
  return cleaned, nil
}

func getCleanedBody(msg string, badWords map[string]bool) string {
  separatedMessage := strings.Split(msg, " ")

  for idx, word := range separatedMessage {
    if badWords[strings.ToLower(word)] == true {
      separatedMessage[idx] = "****"
    } 
  }
  
  filteredMessage := strings.Join(separatedMessage, " ")
  return filteredMessage
}

type Chirp struct {
  ID        uuid.UUID `json:"id"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
  Body      string    `json:"body"`
  UserID    uuid.UUID `json:"user_id"`
}

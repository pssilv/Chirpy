package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenType string

const (
  TokenTypeAcess TokenType = "chirpy-acess"
)

var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")

func HashPassword(password string) (string, error) {
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  if err != nil {
    return "", fmt.Errorf("Failed to hash password: %v", err)
  }

  return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) error {
  err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
  if err != nil {
    return fmt.Errorf("Failed to compare: %v", err)
  }

  return nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
  claims :=  jwt.RegisteredClaims {
    Issuer: string(TokenTypeAcess),
    IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
    ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
    Subject: fmt.Sprintf("%v", userID),
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

  tokenString, err := token.SignedString([]byte(tokenSecret))
  if err != nil {
    return "", fmt.Errorf("Failed  to create JWT: %v", err)
  }

  return tokenString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
  token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(_ *jwt.Token) (interface{}, error) {
    return []byte(tokenSecret), nil
  })
  if err != nil {
    return uuid.Nil, fmt.Errorf("Validation failed: %v", err)
  }

  subject, err := token.Claims.GetSubject()
  if err != nil {
    return uuid.Nil, err
  }

  issuer, err := token.Claims.GetIssuer()
  if err != nil {
    return uuid.Nil, errors.New("Invalid issuer")
  }
  if issuer != string(TokenTypeAcess) {
    return uuid.Nil, errors.New("Invalid issuer")
  }

  userID, err := uuid.Parse(subject)
  if err != nil {
    return uuid.Nil, fmt.Errorf("Failed to get userID: %v", err)
  }

  return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
  authToken := headers.Get("Authorization")
  if authToken == "" {
    return "", ErrNoAuthHeaderIncluded
  }
  splitAuth := strings.Split(authToken, " ")
  if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
    return "", errors.New("malformed Authorization header")
  }

  return splitAuth[1], nil
}

func MakeRefreshToken() (string, error) {
  refreshToken := make([]byte, 32) 
  _, err := rand.Read(refreshToken)
  if err != nil {
    return "", err
  }
  return hex.EncodeToString(refreshToken), nil
}

func GetAPIKey(headers http.Header) (string, error) {
  authApiKey := headers.Get("Authorization")
  if authApiKey == "" {
    return "", ErrNoAuthHeaderIncluded
  }
  splitAuth := strings.Split(authApiKey, " ")
  if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
    return "", errors.New("Malformed Authorization header")
  }

  return splitAuth[1], nil
}

package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
  password1 := "password1231235"
  password2 := "anotherpassword1230213"
  hash1, _ := HashPassword(password1)
  hash2, _ := HashPassword(password2)

  tests := []struct {
    name     string
    password string 
    hash     string 
    wantErr  bool
  }{
    {name: "Correct password", password: password1, hash: hash1, wantErr: false},
    {name: "Incorrect password", password: "Wrongpassword", hash: hash1, wantErr: true},
    {name: "Password with different hash", password: password1, hash: hash2, wantErr: true},
    {name: "Empty password", password: "", hash: hash1, wantErr: true},
    {name: "Invalid hash", password: password1, hash: "a_hash", wantErr: true},
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      err := CheckPasswordHash(tt.password, tt.hash)
      if (err != nil) != tt.wantErr {
        t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func TestMakeValidateJWT(t *testing.T) {
  userID := uuid.New() 
  validToken, _ := MakeJWT(userID, "secret", time.Hour)

  tests := []struct {
    name        string
    tokenString string
    tokenSecret string
    wantUserID uuid.UUID
    wantErr    bool
  }{
    {name: "Valid token", tokenString: validToken, tokenSecret: "secret", wantUserID: userID, wantErr: false},
    {name: "Invalid token", tokenString: "Invalid_token", tokenSecret: "secret", wantUserID: uuid.Nil, wantErr: true},
    {name: "Wrong secret", tokenString: validToken, tokenSecret: "WrongSecret", wantUserID: uuid.Nil, wantErr: true},
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t* testing.T) {
      gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
      if (err != nil) != tt.wantErr {
        t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if gotUserID != tt.wantUserID {
        t.Errorf("ValidateJWT() userID = %v, want %v", userID, tt.wantUserID)
      }

    })
  }
}

func TestGetBearerToken(t *testing.T) {
  tests := []struct {
    name      string
    headers   http.Header
    wantToken string
    wantErr   bool
  }{
    {
      name: "Valid Bearer token", 
      headers: http.Header{
        "Authorization": []string{"Bearer valid_token"},
      },
      wantToken: "valid_token",
      wantErr: false,
    },
    {
      name: "Invalid Bearer token",
      headers: http.Header{
        "Authorization": []string{"InvalidBearer token"},
      },
      wantToken: "",
      wantErr: true,
    },
    {
      name: "Missing Authorization header",
      headers: http.Header{},
      wantToken: "",
      wantErr: true,
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      gotToken, err := GetBearerToken(tt.headers)
      if (err != nil) != tt.wantErr {
        t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if gotToken != tt.wantToken {
        t.Errorf("GetBearerToken() gotToken = %v, want %v", err, tt.wantToken)
      }
    })  
  }
}

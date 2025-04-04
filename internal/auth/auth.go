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
)

type CustomClaims struct {
	jwt.RegisteredClaims
}

func MakeJWT(userId uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userId.String(),
	})
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}
	userId, err := uuid.Parse(subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid userId: %v", err)
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	} else if issuer != "chirpy" {
		return uuid.Nil, errors.New("invalid issuer")
	}

	return userId, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	header := headers.Get("AUTHORIZATION")
	if header == "" {
		return "", errors.New("no Authorization header provided")
	}

	splits := strings.Fields(header)
	if len(splits) < 2 {
		return "", errors.New("invalid Authorization header provided")
	}

	return splits[1], nil
}

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	encoded := hex.EncodeToString(key)

	return encoded, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	header := headers.Get("AUTHORIZATION")
	if header == "" {
		return "", errors.New("no Authorization header provided")
	}

	splits := strings.Fields(header)
	if len(splits) < 2 {
		return "", errors.New("invalid Authorization header provided")
	}

	return splits[1], nil
}

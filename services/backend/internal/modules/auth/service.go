package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	secret []byte
}

func NewService(secret string) *Service { return &Service{secret: []byte(secret)} }

func (s *Service) IssueAccessToken(userID string) (string, error) {
	if userID == "" {
		return "", errors.New("empty user id")
	}
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
}

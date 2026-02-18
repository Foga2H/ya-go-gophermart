package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/Foga2H/ya-go-gophermart/internal/logger"
	"github.com/golang-jwt/jwt/v5"
)

const TokenExp = 3 * time.Hour

var ErrInvalidToken = errors.New("invalid token")
var ErrEmptyUserID = errors.New("empty user ID")

type Claims struct {
	jwt.RegisteredClaims
	UserID string `json:"uid"`
}

type Token struct {
	secret []byte
	logger *logger.Logger
}

func NewToken(secret string, logger *logger.Logger) *Token {
	return &Token{
		secret: []byte(secret),
		logger: logger,
	}
}

func (t *Token) GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: userID,
	})

	return token.SignedString(t.secret)
}

func (t *Token) ParseUserID(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return t.secret, nil
	})

	if err != nil {
		t.logger.Errorf("failed to parse token: %v", err)
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

	if !token.Valid {
		t.logger.Error(ErrInvalidToken)
		return "", ErrInvalidToken
	}

	if claims.UserID == "" {
		t.logger.Error(ErrEmptyUserID)
		return "", ErrEmptyUserID
	}

	return claims.UserID, nil
}

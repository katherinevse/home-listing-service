package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"time"
)

type TokenService struct {
	JWTSecretKey string
	logger       *slog.Logger
}

func NewTokenService(secretKey string, logger *slog.Logger) *TokenService {
	return &TokenService{
		JWTSecretKey: secretKey,
		logger:       logger,
	}
}

type Claims struct {
	UserID   int    `json:"user_id"`
	UserType string `json:"user_type"`
	Email    string `json:"email"` //TODO DEL ME отладка
	jwt.RegisteredClaims
}

func (t *TokenService) GenerateJWT(userID int, userType string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   userID,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(t.JWTSecretKey))
	if err != nil {
		t.logger.Error("Failed to sign JWT token", "error", err)
		return "", fmt.Errorf("failed to sign JWT token: %v", err)
	}
	t.logger.Info("JWT token generated successfully", "userID", userID)
	return signedToken, nil
}

func (t *TokenService) ParseJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	// ParseWithClaims разбирает токен и заполняет структуру claims.
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.JWTSecretKey), nil
	})

	if err != nil {
		t.logger.Error("Error parsing JWT", "error", err)
		return nil, errors.New("invalid token")
	}

	if !token.Valid {
		t.logger.Error("Invalid JWT token")
		return nil, errors.New("invalid token")
	}

	t.logger.Info("JWT token parsed successfully", "userID", claims.UserID)
	return claims, nil
}

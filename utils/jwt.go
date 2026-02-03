package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims is a custom claim structure
type JWTClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims 
}

// GenerateAccessToken to generate access token (30 minutes)
func GenerateAccessToken(UserID, secret string) (string, time.Time, error) {
	expiresAt := time.Now().Add(30 * time.Minute)

	claims := JWTClaims{
		UserID: UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: UserID,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, expiresAt, err
}

// GenerateRefreshToken to generate refresh token (30 days)
func GenerateRefreshToken(userID, secret string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken to validate token and return the claims
func ValidateToken(tokenString string, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
package utils

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/musabgulfam/pumplink-backend/config"
)

func ValidateJWT(tokenStr string) (string, error) {
	cfg := config.Load()
	jwtSecret := []byte(cfg.JWTSecret)
	log.Printf("[JWT] Validating token: %s", tokenStr)

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method if needed:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("[JWT] Unexpected signing method: %v", token.Header["alg"])
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		log.Printf("[JWT] Parse error: %v", err)
		return "", errors.New("invalid token")
	}
	if !token.Valid {
		log.Printf("[JWT] Token is not valid")
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("[JWT] Invalid token claims type: %T", token.Claims)
		return "", errors.New("invalid token claims")
	}

	userIDFloat, ok := claims["sub"].(float64)
	if !ok {
		log.Printf("[JWT] User ID claim missing or wrong type in token claims: %v", claims)
		return "", errors.New("user id claim missing")
	}
	userID := fmt.Sprintf("%.0f", userIDFloat)

	log.Printf("[JWT] Token valid for user: %s", userID)
	return userID, nil
}

package account

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID string `json:"userID"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID string) (string, string, error) {
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)), // token expires in 1 hour
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret
	accessToken, err := token.SignedString([]byte(AppConfig.JWT_SECRET))
	if err != nil {
		return "", "", err
	}
	refreshToken, err := token.SignedString([]byte(AppConfig.REFRESH_JWT_SECRET))
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func VerifyJWT(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(AppConfig.JWT_SECRET), nil
	})
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		log.Println(claims)
		return claims, nil
	} else {
		return nil, err
	}
}

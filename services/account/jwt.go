package account

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key")

// Claims struct (can include custom fields)
type CustomClaims struct {
	UserID string `json:"user_id"`
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
	jwtSecret := AppConfig.JWT_SECRET
	refreshjwtSecret := AppConfig.REFRESH_JWT_SECRET
	// Sign the token with the secret
	accessToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := token.SignedString(refreshjwtSecret)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func VerifyJWT(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

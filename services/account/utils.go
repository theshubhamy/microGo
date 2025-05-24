package account

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func CompareHashPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func isAllowedKey(key string) bool {
	allowedKeys := map[string]struct{}{
		"id":    {},
		"email": {},
		"phone": {},
	}
	_, isExists := allowedKeys[key]
	return isExists
}

func checkPhoneorEmail(emailOrPhone string) (string, error) {
	emailRegex := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	phoneRegex := regexp.MustCompile(`^\+?\d{10,15}$`)

	switch {
	case emailRegex.MatchString(emailOrPhone):
		return "email", nil
	case phoneRegex.MatchString(emailOrPhone):
		return "phone", nil
	default:
		return "", errors.New("input must be a valid email or phone number")
	}
}

func generateFingerprint(ip, userAgent string) string {
	data := ip + userAgent
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

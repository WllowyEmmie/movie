package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwt_secret_key = []byte(os.Getenv("JWT_SECRET_KEY)"))

func GenerateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwt_secret_key)
}

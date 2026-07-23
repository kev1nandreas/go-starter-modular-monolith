package auth

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kev1nandreas/go-rest-api-template/config"
	"golang.org/x/crypto/bcrypt"
)

// Claims struct to be encoded to JWT
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var JwtKey = []byte(config.GetEnvString("JWT_SECRET_KEY", GenerateRandomKey()))

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute).Unix()

	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime,
		Issuer:    username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRandomKey() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic("Failed to generate random key: " + err.Error())
	}

	return base64.StdEncoding.EncodeToString(key)
}

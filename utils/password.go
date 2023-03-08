package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GeneratePassword(password []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate password %s", err)
	}

	return string(hash), nil
}

func ValidatePassword(password []byte, hashPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), password)
}

type CustomClaims struct {
	jwt.StandardClaims
	Email string `json:"email"`
}

func GenerateJwtToken(email string) (string, CustomClaims, error) {
	claims := CustomClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.NewString(),
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
		},
		Email: email,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", CustomClaims{}, err
	}

	return signedToken, claims, nil
}

func ValidateJwtToken(jwtToken string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(jwtToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims type")
	}

	err = claims.Valid()
	if err != nil {
		return nil, err
	}

	return claims, nil
}

package jwtToken

import (
	"errors"
	"iv_project/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTServices interface {
	GenerateToken(userID string, role models.UserRoleType) (string, error)
	DecodeToken(tokenString string) (jwt.MapClaims, error)
}

type jwtService struct {
	secretKey string
	issuer    string
}

func JWTService(secretKey string, issuer string) JWTServices {
	return &jwtService{
		secretKey: secretKey,
		issuer:    issuer,
	}
}

func (j *jwtService) GenerateToken(userID string, role models.UserRoleType) (string, error) {
	claims := jwt.MapClaims{
		"id":   userID,
		"role": role.String(),
		"exp":  time.Now().Add(time.Hour * 48).Unix(), // Token berlaku 48 jam
		"iss":  j.issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j *jwtService) DecodeToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

package utils

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/vtv-us/kahoot-backend/internal/entities"
)

type JwtWrapper struct {
	SecretKey string
	Issuer    string
}

type jwtClaims struct {
	jwt.StandardClaims
	UserID string
	Email  string
}

func (w *JwtWrapper) GenerateToken(user entities.User, expirationHours int64) (signedToken string, err error) {
	expiredAt := time.Now().Local().Add(time.Hour * time.Duration(expirationHours)).Unix()
	claims := &jwtClaims{
		UserID: user.UserID,
		Email:  user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredAt,
			Issuer:    w.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(w.SecretKey))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (w *JwtWrapper) ValidateToken(signedToken string) (claims *jwtClaims, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&jwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(w.SecretKey), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*jwtClaims)

	if !ok {
		return nil, errors.New("couldn't parse claims")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("JWT is expired")
	}

	return claims, nil

}

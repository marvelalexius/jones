package str

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTClaims struct {
	UserID string
	jwt.RegisteredClaims
}

func GenerateJWT(userID string, lifespan time.Time, secret string) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(lifespan),
		},
	}
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := unsignedToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ParseJWT(token string, secret string) (*JWTClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(*JWTClaims); ok && parsedToken.Valid {
		return claims, nil
	}

	return nil, err
}

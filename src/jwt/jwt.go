package jwt

import (
	"errors"
	"fmt"
	"gin_main/config"
	"gin_main/src/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTHelper struct {
	jwtSecret []byte
	ttl       time.Duration
}

type Claims struct {
	Login string
	jwt.RegisteredClaims
}

func NewJWTHelper(cfg *config.Config) *JWTHelper {
	return &JWTHelper{jwtSecret: []byte(cfg.JWT.SecretKey), ttl: cfg.JWT.TTL}
}

func (helper *JWTHelper) ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return helper.jwtSecret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token expired")
		}
		return nil, fmt.Errorf("invalid token")
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

func (helper *JWTHelper) GenerateJWT(user models.User) (string, error) {
	claims := Claims{
		Login: user.Login,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(helper.ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(helper.jwtSecret)
}

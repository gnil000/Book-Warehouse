package jwt

import (
	"errors"
	"fmt"
	"gin_main/config"
	"gin_main/src/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

type JWTHelper struct {
	jwtSecret []byte
	ttl       time.Duration
	log       zerolog.Logger
}

type Claims struct {
	Login string `json:"login"`
	jwt.RegisteredClaims
}

func NewJWTHelper(cfg *config.Config, log zerolog.Logger) *JWTHelper {
	return &JWTHelper{jwtSecret: []byte(cfg.JWT.SecretKey), ttl: cfg.JWT.TTL, log: log}
}

func (helper *JWTHelper) ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			helper.log.Error().Msg("unexpected signing method")
			return nil, fmt.Errorf("unexpected signing method")
		}
		return helper.jwtSecret, nil
	})
	helper.log.Info().Msg(fmt.Sprintf("ТОКЕН пришедший клэимс %v", *token))
	helper.log.Err(err)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token expired")
		}
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(*Claims)
	helper.log.Info().Msg(fmt.Sprintf("Клэимс %v и флаг успеха %t токен.Валид %t", claims, ok, token.Valid))
	if ok && token.Valid {
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
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(helper.jwtSecret)
}

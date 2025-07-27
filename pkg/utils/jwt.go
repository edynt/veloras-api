package utils

import (
	"fmt"
	"time"

	"github.com/edynnt/veloras-api/pkg/global"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type PayloadClaims struct {
	jwt.StandardClaims
}

func GenTokenJWT(payload jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(global.Config.JWT.ApiSecret))
}

func CreateToken(uuidToken string, isRefreshToken bool) (string, error) {
	var (
		expireValue int
		expiration  time.Duration
	)

	if isRefreshToken {
		expireValue = global.Config.JWT.RefreshTokenExpire
		if expireValue <= 0 {
			expireValue = 1
		}
		expiration = time.Duration(expireValue) * 24 * time.Hour // days → duration
	} else {
		expireValue = global.Config.JWT.AccessTokenExpire
		if expireValue <= 0 {
			expireValue = 1
		}
		expiration = time.Duration(expireValue) * time.Hour // hours → duration
	}

	now := time.Now()
	expiresAt := now.Add(expiration)

	return GenTokenJWT(&PayloadClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.New().String(),
			ExpiresAt: expiresAt.Unix(),
			IssuedAt:  now.Unix(),
			Issuer:    "veloras-api",
			Subject:   uuidToken,
		},
	})
}

func ParseJwtTokenSub(token string) (*jwt.StandardClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", jwtToken.Header["alg"])
		}
		return []byte(global.Config.JWT.ApiSecret), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*jwt.StandardClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

func VerifyTokenSubject(token string) (*jwt.StandardClaims, error) {
	claims, err := ParseJwtTokenSub(token)
	if err != nil {
		return nil, err
	}

	if err = claims.Valid(); err != nil {
		return nil, err
	}

	return claims, nil
}

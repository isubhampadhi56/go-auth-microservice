package jwtauth

import (
	"time"

	"github.com/go-auth-microservice/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

type JWT interface {
	CreateToken(jwt.MapClaims) (string, error)
	VerifyToken(string) (jwt.MapClaims, error)
}

var accessTokenHandler JWT
var refreshTokenHandler JWT

func GetAccessTokenHandler() JWT {
	appConfig := config.GetConfig()
	if accessTokenHandler == nil {
		expTime := time.Now().Add(time.Minute * time.Duration(appConfig.GetAccessTokenExpiry())).Unix()
		secret := appConfig.GetAccessTokenSecret()
		accessTokenHandler = InitializeJWTManager(secret, expTime)
	}
	return accessTokenHandler
}
func GetRefreshTokenHandler() JWT {
	appConfig := config.GetConfig()
	if refreshTokenHandler == nil {
		expTime := time.Now().Add(time.Hour * time.Duration(appConfig.GetAccessTokenExpiry())).Unix()
		secret := appConfig.GetRefreshTokenSecret()
		refreshTokenHandler = InitializeJWTManager(secret, expTime)
	}
	return refreshTokenHandler
}

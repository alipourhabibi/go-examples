package services

import (
	"fmt"
	"time"

	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/repo"
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/settings"
	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

type Claims struct {
        UserName string
        jwt.StandardClaims
}

type TokenDetails struct {
	AccessToken string
	RefreshToken string
	AccessTokenUuid string
	RefreshTokenUuid string
	ATExpires int64
	RTExpires int64
}

// Generates Only one JWT based on the given secret and expiration time
func GenerateJWT(username, secret string, expirationTime int64) (string, error) {

        claims := &Claims{
                UserName: username,
                StandardClaims: jwt.StandardClaims{
                        ExpiresAt: expirationTime,
                },
        }

        secretKey := []byte(secret)
        
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        tokenString, err := token.SignedString(secretKey)
        if err == nil {
                return tokenString, nil
        }
        return "", err
}

// Verify JWT based on given secret which can be for accesstoken and refreshtoken
func VerifyJWT(tokenString, secret string) (string, error) {
        claims := &Claims{}
        secretKey := []byte(secret)
        
        token, err := jwt.ParseWithClaims(tokenString, claims, func (token *jwt.Token) (interface{}, error) {
                return secretKey, nil
        })

        if err != nil {
                return "", err
        }

        if token.Valid && claims.ExpiresAt > time.Now().Unix() {
                return claims.UserName, nil
        }

        return "", err
}

// Generate Both access token and refresh token and set them in redis as well
func CreateTokens(username string) (*TokenDetails, error) {
	td := &TokenDetails{}

	accessExp := time.Now().Add(time.Minute * 30).Unix()
	accessToken, _ := GenerateJWT(username, settings.AppSettings.Items.JwtAccess, accessExp)
	td.AccessToken = accessToken
	td.ATExpires = accessExp
	td.RefreshTokenUuid = uuid.NewV4().String()

	refreshExp := time.Now().Add(time.Hour * 30 * 7).Unix()
	refreshToken , _ := GenerateJWT(username, settings.AppSettings.Items.JwtRefresh, refreshExp)
	td.RefreshToken = refreshToken
	td.RTExpires = refreshExp
	td.AccessTokenUuid = td.RefreshTokenUuid + "++" + username

	// save tokens metadata to redis
	at := time.Unix(td.ATExpires, 0)
	rt := time.Unix(td.RTExpires, 0)
	now := time.Now()

	redisClient := repo.GetRedisClient()

	ATCreated, err := redisClient.Set(td.AccessTokenUuid, username, at.Sub(now)).Result()
	if err != nil {
		return nil, err
	}
	RTCreated, err := redisClient.Set(td.RefreshTokenUuid, username, rt.Sub(now)).Result()
	if err != nil {
		return nil, err
	}
	if ATCreated == "0" || RTCreated == "0" {
		return nil, fmt.Errorf("no record inserted")
	}


	return td, nil
}

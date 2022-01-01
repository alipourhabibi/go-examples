package services

import (
	"fmt"
	"time"

	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/repo"
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/settings"
	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

type TokenDetails struct {
	AccessToken string
	RefreshToken string
	AccessTokenUuid string
	RefreshTokenUuid string
	ATExpires int64
	RTExpires int64
}

type Claims struct {
	TD TokenDetails
        jwt.StandardClaims
}

// Generates Only one JWT based on the given secret and expiration time
func GenerateJWT(username string) (*TokenDetails, error) {

	var err error
	td := &TokenDetails{}
	
	accessExp := time.Now().Add(time.Minute * 30).Unix()
	td.ATExpires = accessExp
	td.AccessTokenUuid = uuid.NewV4().String()

	refreshExp := time.Now().Add(time.Hour * 30 * 7).Unix()
	td.RTExpires = refreshExp
	td.RefreshTokenUuid = td.AccessTokenUuid + "++" + username

        secretKeyAcess := []byte(settings.AppSettings.Items.JwtAccess)
	secretKeyRefresh := []byte(settings.AppSettings.Items.JwtRefresh)

	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = td.AccessTokenUuid
	atClaims["username"] = username
	atClaims["exp"] = td.ATExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString(secretKeyAcess)
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshTokenUuid
	rtClaims["username"] = username
	rtClaims["exp"] = td.RTExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString(secretKeyRefresh)
	if err != nil {
		return nil, err
	}

	return td, nil
        
}

// Verify JWT based on given secret which can be for accesstoken and refreshtoken
func VerifyJWT(tokenString, secret string) (*jwt.Token, error) {
	claims := jwt.MapClaims{}
        secretKey := []byte(secret)
        
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
                return secretKey, nil
        })

        if err != nil {
                return nil, err
        }

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Internal Server Error on verifying jwt")
	}
	exp, _ := claims["exp"].(float64)
        if token.Valid && time.Unix(int64(exp), 0).Unix() > time.Now().Unix() {
                return token, nil
        }
        return nil, err
}

// Generate Both access token and refresh token and set them in redis as well
func CreateTokensAndMetaData(username string) (*TokenDetails, error) {

	td, err := GenerateJWT(username)

	if err != nil {
		return nil, err
	}

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

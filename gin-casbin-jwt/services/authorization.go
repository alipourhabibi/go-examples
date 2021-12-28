package services

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
        UserName string
        jwt.StandardClaims
}

func GenerateJWT(username, secret string) (string, error) {
        // for 30 days
        expirationTime := time.Now().Add(time.Hour * 24 * 30).Unix()
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
	

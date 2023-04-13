package pkg

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Auth Jwt auth
type Auth struct{}

// NewAuth New Auth
func NewAuth() *Auth {
	return &Auth{}
}

// Claims Custom claims
type Claims struct {
	jwt.MapClaims
	Username string `json:"username"`
}

// GenerateToken Generate jwt token
func (a *Auth) GenerateToken(username string, duration int64) (string, error) {

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	keyFile, err := os.Open(dir + "/pkg/prikey.pem")
	if err != nil {
		panic(err)
	}
	var keyByte []byte
	_, err = keyFile.Read(keyByte)
	if err != nil {
		panic(err)
	}

	fmt.Println(keyByte)

	priKey, err := rsa.GenerateKey(keyFile, 2048)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(priKey)

	tokenObj := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":      "rgrpc",
		"aud":      "rgrpc",
		"nbf":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Second * 10).Unix(),
		"sub":      "user",
		"username": username,
	})
	token, err := tokenObj.SignedString(priKey)
	iss, err := tokenObj.Claims.GetIssuer()
	if err != nil {
		panic(err)
	}
	fmt.Println("iss is" + iss)
	//token, err := tokenObj.SignedString([]byte("$2a$04$cCglP2MoMEL8wjP8JJm0..wl4CzfZwTx7ZZEoVsvb.OnGPGaIxnhS"))
	if err != nil {
		return "", err
	}
	return token, nil

}

// ValidateAuth Validate jwt auth
func (a *Auth) ValidateAuth(authToken string) (string, error) {

	claims := &Claims{}
	jwtToken, err := jwt.ParseWithClaims(authToken, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Header["alg"] != "HS256" {
			return "", errors.New("InvalidAlgorithm")
		}
		return []byte("$2a$04$cCglP2MoMEL8wjP8JJm0..wl4CzfZwTx7ZZEoVsvb.OnGPGaIxnhS"), nil
	})
	fmt.Println(jwtToken.Claims.GetIssuer())

	if err != nil {
		return "", err
	}

	if !jwtToken.Valid {
		return "", err
	}

	return claims.Username, nil

}

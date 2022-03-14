package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

const subject = "ipashare"

var jwtSecret = uuid.New().String()

type LoginClaims struct {
	jwt.StandardClaims
	Username string
}

func GenerateJWT(username string, minute ...int) (string, error) {
	nowTime := time.Now()
	var expireTime time.Time
	if len(minute) > 0 {
		expireTime = nowTime.Add(time.Duration(minute[0]) * time.Minute)
	} else {
		expireTime = nowTime.Add(8 * time.Hour)
	}
	loginClaims := LoginClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
			Subject:   subject,
		},
		Username: username,
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, loginClaims)
	return tokenClaims.SignedString([]byte(jwtSecret))
}

func ParseJWT(token string) (*LoginClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &LoginClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
	if err != nil {
		return nil, err
	}
	if jwtToken != nil {
		if claims, ok := jwtToken.Claims.(*LoginClaims); ok && jwtToken.Valid {
			return claims, nil
		}
	}
	return nil, errors.New("parse jwt fail")
}

func ChangeJwtSecret() {
	jwtSecret = uuid.New().String()
}

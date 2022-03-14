package cache

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"ipashare/pkg/auth"
)

const (
	_errCode   = "errcode"
	_jwtClaims = "jwtClaims"
)

func GetCode(c *gin.Context) string {
	if value, ok := c.Get(_errCode); ok {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

func SetCode(c *gin.Context, code int, msg string) {
	c.Set(_errCode, fmt.Sprintf("code = %d msg = %s", code, msg))
}

func GetJwtClaims(c *gin.Context) *auth.LoginClaims {
	if value, ok := c.Get(_jwtClaims); ok {
		if claims, ok := value.(*auth.LoginClaims); ok {
			return claims
		}
	}
	return nil
}

func SetJwtClaims(c *gin.Context, claims *auth.LoginClaims) {
	c.Set(_jwtClaims, claims)
}

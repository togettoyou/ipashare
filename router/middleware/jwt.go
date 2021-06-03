package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
	. "super-signature/handler"
	"super-signature/util/errno"
	"super-signature/util/tools"
)

var jwtClaims = "jwtClaims"

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		g := Gin{Ctx: c}
		auth := c.Request.Header.Get("Authorization")
		if len(auth) == 0 {
			g.SendNoDataResponse(errno.ErrNotLogin)
			c.Abort()
			return
		}
		if len(strings.Fields(auth)) > 1 {
			auth = strings.Fields(auth)[1]
		}
		// 校验token
		claims, err := tools.ParseJWT(auth)
		if err != nil {
			if validationError, ok := err.(*jwt.ValidationError); ok {
				switch validationError.Errors {
				case jwt.ValidationErrorExpired:
					g.SendNoDataResponse(errno.ErrTokenExpired)
					break
				default:
					g.SendNoDataResponse(errno.ErrTokenInvalid)
					break
				}
				c.Abort()
				return
			}
			g.SendNoDataResponse(errno.New(errno.ErrTokenFailure, err).Add(err.Error()))
			c.Abort()
			return
		}
		c.Set(jwtClaims, claims)
		c.Next()
	}
}

func GetJWTClaims(c *gin.Context) *tools.Claims {
	if claims, ok := c.Get(jwtClaims); ok {
		if cs, ok := claims.(*tools.Claims); ok {
			return cs
		}
	}
	return nil
}

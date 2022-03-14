package middleware

import (
	"ipashare/internal/api"
	"ipashare/internal/server/middleware/cache"
	"ipashare/pkg/auth"
	"ipashare/pkg/e"
	"ipashare/pkg/log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtApi := api.Base{
			Log: log.New("JWT").L(),
		}
		jwtApi.MakeContext(c)
		authorization := c.Request.Header.Get("Authorization")
		if len(authorization) == 0 {
			jwtApi.Resp(http.StatusForbidden, e.ErrNotLogin, false)
			return
		}
		if len(strings.Fields(authorization)) > 1 {
			authorization = strings.Fields(authorization)[1]
		}
		claims, err := auth.ParseJWT(authorization)
		if err != nil {
			if validationError, ok := err.(*jwt.ValidationError); ok {
				switch validationError.Errors {
				case jwt.ValidationErrorExpired:
					jwtApi.Resp(http.StatusForbidden, e.ErrTokenExpired, false)
					return
				default:
					jwtApi.Resp(http.StatusForbidden, e.ErrTokenInvalid, false)
					return
				}
			}
			jwtApi.Resp(http.StatusForbidden, e.ErrTokenFailure, false)
			return
		}
		cache.SetJwtClaims(c, claims)
		c.Next()
	}
}

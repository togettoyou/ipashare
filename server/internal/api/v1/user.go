package v1

import (
	"ipashare/internal/api"
	"ipashare/internal/model/req"
	"ipashare/internal/server/middleware/cache"
	"ipashare/internal/svc"

	"github.com/gin-gonic/gin"
)

type User struct {
	api.Base
}

// Login
// @Tags User
// @Summary 登录
// @Produce json
// @Param data body req.UserPW true "登录信息"
// @Success 200 {object} api.Response
// @Router /api/v1/user/login [post]
func (u User) Login(c *gin.Context) {
	var (
		body    req.UserPW
		userSvc svc.User
	)
	if !u.MakeContext(c).MakeService(&userSvc.Service).ParseJSON(&body) {
		return
	}
	info, err := userSvc.Login(body.Username, body.Password)
	if u.HasErr(err) {
		return
	}
	u.OK(info)
}

// ChangePW
// @Tags User
// @Summary 修改密码
// @Security ApiKeyAuth
// @Produce json
// @Param data body req.UserPW true "登录信息"
// @Success 200 {object} api.Response
// @Router /api/v1/user/changepw [post]
func (u User) ChangePW(c *gin.Context) {
	var (
		body    req.UserPW
		userSvc svc.User
	)
	if !u.MakeContext(c).MakeService(&userSvc.Service).ParseJSON(&body) {
		return
	}
	err := userSvc.ChangePW(cache.GetJwtClaims(c).Username, body.Username, body.Password)
	if u.HasErr(err) {
		return
	}
	u.OK()
}

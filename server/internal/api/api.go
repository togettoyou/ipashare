package api

import (
	"fmt"
	"net/http"

	"ipashare/internal/model"
	"ipashare/internal/server/middleware/cache"
	"ipashare/internal/svc"
	"ipashare/pkg/e"
	logpkg "ipashare/pkg/log"
	"ipashare/pkg/validatorer"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// Base 每个api结构体都需要内嵌该结构体
type Base struct {
	log   *zap.Logger
	store *model.Store
	c     *gin.Context

	Log *zap.Logger
}

func New(store *model.Store, log *zap.Logger) Base {
	if log == nil {
		log = logpkg.New("").L()
	}
	return Base{
		log:   log,
		store: store,
		Log:   log.Named("api"),
	}
}

// MakeContext 初始化http上下文
func (b *Base) MakeContext(c *gin.Context) *Base {
	b.c = c
	return b
}

// MakeService 初始化业务对象
func (b *Base) MakeService(svc ...*svc.Service) *Base {
	for _, service := range svc {
		if service != nil {
			service.New(b.store, b.log)
		}
	}
	return b
}

func (b *Base) Named(name string) *Base {
	if b.log == nil {
		b.log = logpkg.New("").L()
	}
	b.log = b.log.Named(name)
	b.Log = b.log.Named("api")
	return b
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// OK 成功响应
func (b *Base) OK(arg ...interface{}) {
	resp := Response{Code: 0, Msg: "OK", Data: map[string]interface{}{}}
	if len(arg) > 0 {
		resp.Data = arg[0]
	}
	cache.SetCode(b.c, resp.Code, resp.Msg)
	b.c.AbortWithStatusJSON(http.StatusOK, resp)
}

// Resp 自定义响应
func (b *Base) Resp(httpCode int, err error, log bool, arg ...interface{}) {
	code, msg := e.DecodeErr(err)
	resp := Response{Code: code, Msg: msg, Data: map[string]interface{}{}}
	if len(arg) > 0 {
		resp.Data = arg[0]
	}
	cache.SetCode(b.c, code, msg)
	b.c.AbortWithStatusJSON(httpCode, resp)
	if log && err != nil {
		b.log.Error(fmt.Sprintf("%+v", err))
	}
}

// HasErr 判断业务错误
func (b *Base) HasErr(err error) bool {
	if err != nil {
		b.Resp(http.StatusInternalServerError, err, true)
		return true
	}
	return false
}

// Parse 解析参数并校验，自动解析方式
// obj 要解析的结构体实例地址
// bind 解析类型
// hideDetails 可选择隐藏参数校验详细信息
func (b *Base) Parse(obj interface{}, hideDetails ...bool) bool {
	if err := b.c.ShouldBind(obj); err != nil {
		return b.validatorData(err, len(hideDetails) > 0 && hideDetails[0])
	}
	return true
}

// ParseWith 解析参数并校验，自定义解析方式
// obj 要解析的结构体实例地址
// bind 解析类型
// hideDetails 可选择隐藏参数校验详细信息
func (b *Base) ParseWith(obj interface{}, bind binding.Binding, hideDetails ...bool) bool {
	if err := b.c.ShouldBindWith(obj, bind); err != nil {
		return b.validatorData(err, len(hideDetails) > 0 && hideDetails[0])
	}
	return true
}

// ParseUri 解析参数并校验，uri方式
// obj 要解析的结构体实例地址
// hideDetails 可选择隐藏参数校验详细信息
func (b *Base) ParseUri(obj interface{}, hideDetails ...bool) bool {
	if err := b.c.ShouldBindUri(obj); err != nil {
		return b.validatorData(err, len(hideDetails) > 0 && hideDetails[0])
	}
	return true
}

// ParseQuery 解析参数并校验，query方式
// obj 要解析的结构体实例地址
// hideDetails 可选择隐藏参数校验详细信息
func (b *Base) ParseQuery(obj interface{}, hideDetails ...bool) bool {
	if err := b.c.ShouldBindQuery(obj); err != nil {
		return b.validatorData(err, len(hideDetails) > 0 && hideDetails[0])
	}
	return true
}

// ParseJSON 解析参数并校验，json方式
// obj 要解析的结构体实例地址
// hideDetails 可选择隐藏参数校验详细信息
func (b *Base) ParseJSON(obj interface{}, hideDetails ...bool) bool {
	if err := b.c.ShouldBindJSON(obj); err != nil {
		return b.validatorData(err, len(hideDetails) > 0 && hideDetails[0])
	}
	return true
}

// ParseForm 解析参数并校验，form方式
// obj 要解析的结构体实例地址
// hideDetails 可选择隐藏参数校验详细信息
func (b *Base) ParseForm(obj interface{}, hideDetails ...bool) bool {
	if err := b.c.ShouldBindWith(obj, binding.Form); err != nil {
		return b.validatorData(err, len(hideDetails) > 0 && hideDetails[0])
	}
	return true
}

// validatorData 参数校验，默认使用msg显示错误，可更换为data
// hideDetails 可选择隐藏参数校验详细信息
func (b *Base) validatorData(err error, hideDetails bool) bool {
	if vErr, ok := err.(validator.ValidationErrors); ok {
		if !hideDetails {
			b.Resp(
				http.StatusBadRequest,
				e.New(e.ErrValidation, vErr).Add(validatorer.TranslateErrMsg(vErr)),
				false,
			)
		} else {
			b.Resp(http.StatusBadRequest, e.ErrValidation, false)
		}
		return false
	}
	b.Resp(http.StatusInternalServerError, e.BindError, false)
	return false
}

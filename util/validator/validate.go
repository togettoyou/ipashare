package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"go.uber.org/zap"
	"reflect"
	"strings"
)

var (
	v     *validator.Validate
	trans ut.Translator
)

func V() *validator.Validate {
	return v
}

func E(err error) (string, bool) {
	if err != nil {
		if errors, ok := err.(validator.ValidationErrors); ok {
			return TranslateErrMsg(errors), false
		}
		return err.Error(), false
	}
	return "", true
}

// Setup 定制gin内置的validator
func Setup() {
	var ok bool
	if v, ok = binding.Validator.Engine().(*validator.Validate); ok {
		zhT := zh.New()              // 中文翻译器
		enT := en.New()              // 英文翻译器
		uni := ut.New(enT, zhT, enT) // 支持中英文，不支持时选择回滚英文
		var transOk bool
		trans, transOk = uni.GetTranslator("zh")
		if !transOk {
			zap.L().Error("translator is not ok")
			return
		}
		// 验证器注册翻译器
		err := zhTranslations.RegisterDefaultTranslations(v, trans)
		if err != nil {
			zap.L().Error(err.Error())
			return
		}
		registerTagNameFunc()
		registerValidationTranslation([]validationTranslation{
			{
				tag: "checkUsername",
				Fun: checkUsername,
				msg: "{0}必须是由字母开头的4-16位字母和数字组成的字符串",
			},
		})
	}
}

// TranslateErrMsg 以msg方式翻译错误消息
func TranslateErrMsg(errs validator.ValidationErrors) string {
	var errList []string
	for _, e := range errs {
		errList = append(errList, e.Translate(trans))
	}
	return strings.Join(errList, "|")
}

// TranslateErrData 以data方式翻译错误消息
func TranslateErrData(errs validator.ValidationErrors) map[string]string {
	return removeTopStruct(errs.Translate(trans))
}

// registerTagNameFunc 获取json标签，作为字段名称
func registerTagNameFunc() {
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

type validationTranslation struct {
	tag string
	Fun validator.Func
	msg string
}

// registerValidationTranslation 注册验证方法并翻译
func registerValidationTranslation(vs []validationTranslation) {
	for _, v := range vs {
		registerValidation(v.tag, v.Fun)
		registerTranslation(v.tag, registerTranslator(v.tag, v.msg))
	}
}

// registerValidation 自定义验证方法
func registerValidation(tag string, fun validator.Func) {
	if err := v.RegisterValidation(tag, fun); err != nil {
		zap.L().Error(err.Error())
	}
}

// registerTranslation 根据自定义的标记注册翻译
func registerTranslation(tag string, registerFn validator.RegisterTranslationsFunc) {
	if err := v.RegisterTranslation(
		tag,
		trans,
		registerFn,
		translate); err != nil {
		zap.L().Error(err.Error())
	}
}

// registerTranslator 为自定义字段添加翻译功能
func registerTranslator(tag string, msg string) validator.RegisterTranslationsFunc {
	return func(trans ut.Translator) error {
		if err := trans.Add(tag, msg, false); err != nil {
			return err
		}
		return nil
	}
}

// translate 自定义字段的翻译方法
func translate(trans ut.Translator, fe validator.FieldError) string {
	msg, err := trans.T(fe.Tag(), fe.Field())
	if err != nil {
		return fe.(error).Error()
	}
	return msg
}

// removeTopStruct 去除字段名中的结构体名称标识
// refer from:https://github.com/go-playground/validator/issues/633#issuecomment-654382345
func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}

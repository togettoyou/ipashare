package validatorer

import (
	"errors"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	validate *validator.Validate
	trans    ut.Translator

	transErr = errors.New("translator is not ok")
)

func V() *validator.Validate {
	return validate
}

func E(err error) (string, bool) {
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return TranslateErrMsg(validationErrors), false
		}
		return err.Error(), false
	}
	return "", true
}

// Setup 定制gin内置的validator
func Setup() {
	var ok bool
	if validate, ok = binding.Validator.Engine().(*validator.Validate); ok {
		zhT := zh.New()              // 中文翻译器
		enT := en.New()              // 英文翻译器
		uni := ut.New(enT, zhT, enT) // 支持中英文，不支持时选择回滚英文
		var transOk bool
		trans, transOk = uni.GetTranslator("zh")
		if !transOk {
			panic(transErr)
		}
		// 验证器注册翻译器
		err := zhTranslations.RegisterDefaultTranslations(validate, trans)
		if err != nil {
			panic(err)
		}
		registerTagNameFunc()
		registerValidationTranslation(vt)
	} else {
		panic("validator setup fail")
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
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
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
func registerValidationTranslation(vt []validationTranslation) {
	for _, v := range vt {
		if err := validate.RegisterValidation(v.tag, v.Fun); err != nil {
			panic(err)
		}
		if err := validate.RegisterTranslation(
			v.tag,
			trans,
			registerTranslator(v.tag, v.msg),
			translate); err != nil {
			panic(err)
		}
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

package validatorer

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	vt = []validationTranslation{
		{
			tag: "checkUsername",
			Fun: checkUsername,
			msg: "{0}必须是由字母开头的4-16位字母和数字组成的字符串",
		},
		// 可自行扩展...
	}
)

// checkUsername 检查用户名正则
func checkUsername(fl validator.FieldLevel) bool {
	if ok, _ := regexp.MatchString(`^[a-zA-Z]{1}[a-zA-Z0-9]{3,15}$`, fl.Field().String()); !ok {
		return false
	}
	return true
}

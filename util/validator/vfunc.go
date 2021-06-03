package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

// checkUsername 检查用户名正则
func checkUsername(fl validator.FieldLevel) bool {
	if ok, _ := regexp.MatchString(`^[a-zA-Z]{1}[a-zA-Z0-9]{3,15}$`, fl.Field().String()); !ok {
		return false
	}
	return true
}

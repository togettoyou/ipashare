package e

import (
	"fmt"
	"gorm.io/gorm"

	"github.com/pkg/errors"
)

type errno struct {
	code int
	msg  string
}

func (e errno) Error() string {
	return fmt.Sprintf("code = %d msg = %s", e.code, e.msg)
}

type codeError struct {
	code int
	msg  string
	err  error
}

func (e *codeError) Add(msg string) error {
	e.msg += ": " + msg
	return e
}

func (e *codeError) Addf(format string, args ...interface{}) error {
	e.msg += fmt.Sprintf(format, args...)
	return e
}

func (e *codeError) Error() string {
	return fmt.Sprintf("code = %d msg = %s err = %s", e.code, e.msg, e.err)
}

// New 新建错误
func New(errno *errno, err error) *codeError {
	return &codeError{code: errno.code, msg: errno.msg, err: err}
}

// NewWithStack 新建错误，附加调用栈
func NewWithStack(errno *errno, err error) error {
	return errors.WithStack(&codeError{code: errno.code, msg: errno.msg, err: err})
}

// Wrap 包装错误信息，附加调用栈
// 第二个参数只能是 string，也可以不传，大部分情况不用传
func Wrap(err error, args ...string) error {
	if len(args) >= 1 {
		return errors.Wrap(err, args[0])
	}

	return errors.Wrap(err, "")
}

// Cause 获取原始错误对象
func Cause(err error) error {
	return errors.Cause(err)
}

// DecodeErr 提取错误码，返回 code 和 msg
func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.code, OK.msg
	}
	err = Cause(err)
	switch typed := err.(type) {
	case *errno:
		return typed.code, typed.msg
	case *codeError:
		if typed.err == gorm.ErrRecordNotFound {
			return DBRecordNotFoundError.code, DBRecordNotFoundError.msg
		}
		return typed.code, typed.msg
	}
	return ServerError.code, err.Error()
}

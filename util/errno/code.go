package errno

// 例: 10001
// 第一位的1表示服务级错误 (1 为系统级错误；2 为普通错误，通常是由用户非法操作引起的)
// 第二位至第三位的00表示服务模块代码
// 最后两位01表示具体错误代码

var (
	OK                  = &Errno{Code: 0, Message: "成功"}
	InternalServerError = &Errno{Code: 10001, Message: "服务器异常"}

	ErrValidation = &Errno{Code: 20001, Message: "参数校验失败"}
	ErrBind       = &Errno{Code: 20002, Message: "参数绑定异常"}
	ErrUnknown    = &Errno{Code: 20003, Message: "未知错误"}

	ErrUploadP8  = &Errno{Code: 20201, Message: "请上传p8文件类型"}
	ErrUploadIPA = &Errno{Code: 20202, Message: "请上传IPA文件类型"}
	ErrUploadIss = &Errno{Code: 20203, Message: "iss 不能为空"}
	ErrUploadKid = &Errno{Code: 20204, Message: "kid 不能为空"}

	ErrHaveAppleAccount   = &Errno{Code: 20301, Message: "开发者账号已存在,不能重复添加"}
	ErrNotAppleAccount    = &Errno{Code: 20302, Message: "开发者账号不存在"}
	ErrDeviceInsufficient = &Errno{Code: 20303, Message: "可用设备已不足"}
	ErrNotFile            = &Errno{Code: 20304, Message: "文件不存在"}
	ErrNotIPA             = &Errno{Code: 20305, Message: "IPA包不存在"}
	ErrSignIPA            = &Errno{Code: 20305, Message: "等待签名失败"}
)

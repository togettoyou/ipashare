package e

// 例: 10001
// 第一位的1表示服务级错误 (1 为系统级错误；2 为普通错误，通常是由用户非法操作引起的)
// 第二位至第三位的00表示服务模块代码
// 最后两位01表示具体错误代码

var (
	OK          = &errno{code: 0, msg: "成功"}
	ServerError = &errno{code: 10001, msg: "服务器异常"}
	DBError     = &errno{code: 10002, msg: "DB异常"}
	BindError   = &errno{code: 10003, msg: "参数绑定异常"}

	ErrValidation   = &errno{code: 20001, msg: "参数校验失败"}
	ErrUploadFormat = &errno{code: 20002, msg: "文件上传格式不合法"}

	ErrNotLogin = &errno{code: 20101, msg: "请登录"}
	// 可自行扩展...

	ErrIPAParser = &errno{code: 20201, msg: "IPA 解析异常"}
	ErrIPAIcon   = &errno{code: 20202, msg: "IPA Icon 提取异常"}

	ErrIssExist           = &errno{code: 20301, msg: "Iss 已存在"}
	ErrIssAdd             = &errno{code: 20302, msg: "账号添加失败"}
	ErrAppstoreAPI        = &errno{code: 20303, msg: "App Store Connect API 连接失败"}
	ErrDeviceInsufficient = &errno{code: 20304, msg: "账号可绑设备已不足"}
)

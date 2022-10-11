package e

// 例: 10001
// 第一位的1表示服务级错误 (1 为系统级错误；2 为普通错误，通常是由用户非法操作引起的)
// 第二位至第三位的00表示服务模块代码
// 最后两位01表示具体错误代码

var (
	OK                    = &errno{code: 0, msg: "成功"}
	ServerError           = &errno{code: 10001, msg: "服务器异常"}
	DBError               = &errno{code: 10002, msg: "DB异常"}
	DBRecordNotFoundError = &errno{code: 10003, msg: "找不到对应记录"}
	BindError             = &errno{code: 10004, msg: "参数绑定异常"}

	ErrValidation   = &errno{code: 20001, msg: "参数校验失败"}
	ErrUploadFormat = &errno{code: 20002, msg: "文件上传格式不合法"}

	ErrNotLogin     = &errno{code: 20101, msg: "请登录"}
	ErrPassword     = &errno{code: 20102, msg: "密码错误"}
	ErrTokenGen     = &errno{code: 20103, msg: "令牌生成失败"}
	ErrTokenExpired = &errno{code: 20104, msg: "令牌已过期"}
	ErrTokenInvalid = &errno{code: 20105, msg: "令牌无效"}
	ErrTokenFailure = &errno{code: 20106, msg: "令牌验证失败"}
	// 可自行扩展...

	ErrIPAParser = &errno{code: 20201, msg: "IPA 解析异常"}
	ErrIPAIcon   = &errno{code: 20202, msg: "IPA Icon 提取异常"}

	ErrIssExist             = &errno{code: 20301, msg: "Iss 已存在"}
	ErrIssAdd               = &errno{code: 20302, msg: "账号添加失败"}
	ErrAppstoreAPI          = &errno{code: 20303, msg: "App Store Connect API 连接异常，详细信息请查看日志"}
	ErrCertificateNotEnough = &errno{code: 20304, msg: "iOS Development certificate 数量已满，请删除一个"}
	ErrSign                 = &errno{code: 20305, msg: "打包失败，详细信息请查看日志"}

	ErrOSSVerify = &errno{code: 20401, msg: "OSS校验失败，请检查配置，详细信息请查看日志"}
	ErrOSSEnable = &errno{code: 20402, msg: "OSS未启用"}
)

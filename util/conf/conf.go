package conf

import (
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"super-signature/util/tools"
	"sync"
)

type config struct {
	IPASign   sync.Map
	Model     string
	ApplePath applePath
}

type applePath struct {
	URL                   string
	AppleAccountPath      string
	UploadPath            string
	TemporaryDownloadPath string
}

type csr struct {
	Key     string
	KeyPath string
	Csr     string
	CsrPath string
}

var (
	Config        config
	CSRSetting    csr
	defaultIosCsr = `-----BEGIN CERTIFICATE REQUEST-----
MIICqTCCAZECAQAwZDELMAkGA1UEBhMCQ04xEjAQBgNVBAgTCUd1YW5nRG9uZzES
MBAGA1UEBxMJR3Vhbmd6aG91MQ8wDQYDVQQKDAbohb7orq8xCzAJBgNVBAsTAkdv
MQ8wDQYDVQQDEwZxcS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIB
AQCrVNaLUC0/lRkcE3D+nxhaX9lSRTjFbNuaDaoo4dOuSx9Z664HcdsSoD4ZldpH
0QcYXin4r2nxPTitasQq7fb8sEHGCt5TbsIxUKa81NNqToEPL0kCqM+1wVMjqdoU
hsm/NmvlROkimwWN/r7NGAn2F3ocT0PhbHDSjUorPRUHNMGPNbgeHqgO9wv2Aw4V
1Rxx4QCMQ5jg/usxY3VS6/gUHsX+IiejMyEBLvJfqMKPsZicWNMDbMML2sQdZfha
IytwZ1pJ3+nMvliMLWEFMQMpy2SIjbGdGW39i6hjMPT0lq7cy31wUqkpWlHQDWhU
GOpsLir/8er1Asj4leFsILZHAgMBAAGgADANBgkqhkiG9w0BAQsFAAOCAQEAJokv
D2f449Lycw+grHjbj3jaSPzceG12V/NxyaYj70TtEXw+zwk2cg5go+MEWWTMHCfe
S0wwlz2aITmNVGzVs7mmWPnhgZ4QGqpWNXDvOqjAWDgQX6rIL4Secl0GhhqOg+tM
OntqKEKbLennS/xWUqPpahEc+DkJ5ud9vBYpfXr3quCDMbVtBoKYLZ0PBoaY8Wr/
y3CIXmmdw1rvCAMlAtasCEbU8oGjX4OBBEeV8yfB/c8wf07zCRBj2YjhYtyghsI+
mDZwjmwJcsYlhdLnBWgbc4979Ow+zMlhSExJ1UCIm3EmlZiBBGf3fMPwXdS8k/YU
O+xHc2lD0LYH6BoyJA==
-----END CERTIFICATE REQUEST-----
`
	defaultIosKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAq1TWi1AtP5UZHBNw/p8YWl/ZUkU4xWzbmg2qKOHTrksfWeuu
B3HbEqA+GZXaR9EHGF4p+K9p8T04rWrEKu32/LBBxgreU27CMVCmvNTTak6BDy9J
AqjPtcFTI6naFIbJvzZr5UTpIpsFjf6+zRgJ9hd6HE9D4Wxw0o1KKz0VBzTBjzW4
Hh6oDvcL9gMOFdUcceEAjEOY4P7rMWN1Uuv4FB7F/iInozMhAS7yX6jCj7GYnFjT
A2zDC9rEHWX4WiMrcGdaSd/pzL5YjC1hBTEDKctkiI2xnRlt/YuoYzD09Jau3Mt9
cFKpKVpR0A1oVBjqbC4q//Hq9QLI+JXhbCC2RwIDAQABAoIBAC2z4c8jwg417Y7J
uNiTA+IHs2b4xB4V3baIcp2ZL+hMbb6E2dVuj6u2Rxp4GNQTdDsR00xdLnuFgzrv
QgjZlYruUX1MpOXIo9CX3QJ+Gy8+ZbrxOB6XfWDUgyL+SggztFlnYPy1lyL+C0tH
awo2oWGd/ZrToh3d2XKw8dn630MeZYgtNx7Ns41yx3nEUsbcx5/57/3hGo0XVTbw
vBmXMTvOXjsGS9n3J+g8MoKrZDS4WZvk5XCO6TZqjIpdLTXTRU9aKNusggFVxL3I
hcJJYPaUEc5Wv67Zw/UgJby+/oa4HikvZYHPoJoj1NJfjX4z/nKmS8JBwOKWKcNv
VsGGYXECgYEA4NBRV5TyucYEYuVI2OSjrozEe37NcUgt/Bp+pHI11VOc5SVPWRk7
2OyHyJtuUgKhNhn7BvJLa3aIR5PjsgrvvcJXLbRNmuXhJ+rmyKV6XPaUE4583rBz
q+y2TNAqLkK6AyTVexEADb3Zs8m2P9PAZus8ImXTcDJ7wo1S8LApMq0CgYEAwxk7
fNs49VsK5ZCZOOxu3eheYyP/J3oYcu/LyhQP+h2tDBV4PNw4GJ1sSOpcmuYW0rwG
j2RFi2b+RYW0p6zEREI0uLntIfo56Y5s8F7J9eQvFLoYHbreAc+jcIfqd00at7UR
aXmipKYZkhKTueqMhLHuoVs3l15Zge3lVeE7n0MCgYEAhT6S54Dtd/QIR4Ez8vFY
nizqi3N1Wm34a1JcuyTCCWUcOagqZlmRYhmWxOxyr4LFf+ZYJR7YWqIPVbUuoCjh
PSwBNaKG2IblMx6DmGqToqO20fyCwA3/EOgkiFRcm7yKuTBMoztJN9vNO0UTkrLz
d0x3AMAvWHFjbsUKYoNWd6UCgYBcHvDw3o6Bg9CcXu+KalFbJJpU061qFYOv2bxv
GZQFtLd+CjExA4bVUJfNwiOh3F0QeknasS0JSsNrTlvkBHbUCDLeuqPWtFFeD6su
wIM6QNlePxSeDRtA3as9ul+in1yrO6sSE6YReoB+cZkhYzegfGfB9tFD/v/iktMD
/wrpAwKBgCmALm+45gTO3ZicXIa7tluJ3+SRJLrzcJUCrBXnnVFrBESJM6EXLsbJ
r+IOMYXHzUNgEOO9L0w7g8MvQSwQHaPYHOwLDn+7uNa3fQqzjLkISr78yIcIGUma
34vQzNIk28M1WEFvc1KOyetHKfLS1Zpa+t9K15fU+GqK0icVs33K
-----END RSA PRIVATE KEY-----
`
)

// Setup 读取配置文件设置
func Setup(url, mode, iosCsr, iosKey string) {
	Config = config{
		IPASign: sync.Map{},
		Model:   mode,
		ApplePath: applePath{
			URL:                   url,
			AppleAccountPath:      "./ios/appleAccount/",
			UploadPath:            "./ios/upload/",
			TemporaryDownloadPath: "./ios/temporaryDownload/",
		}}
	CSRSetting = csr{
		KeyPath: iosKey,
		CsrPath: iosCsr,
	}
	setConfig()
}

// setConfig 构造配置文件到Config结构体上
func setConfig() {
	createPath(Config.ApplePath.AppleAccountPath)
	createPath(Config.ApplePath.UploadPath)
	createPath(Config.ApplePath.TemporaryDownloadPath)
	if CSRSetting.CsrPath != "" && CSRSetting.KeyPath != "" {
		keyData, err := ioutil.ReadFile(CSRSetting.KeyPath)
		if err != nil {
			panic(err)
		}
		CSRSetting.Key = string(keyData)
		csrData, err := ioutil.ReadFile(CSRSetting.CsrPath)
		if err != nil {
			panic(err)
		}
		CSRSetting.Csr = string(csrData)
	} else {
		CSRSetting.Csr = defaultIosCsr
		CSRSetting.Key = defaultIosKey
		CSRSetting.CsrPath = "conf/ios.csr"
		CSRSetting.KeyPath = "conf/ios.key"
		err := tools.CreateFile(defaultIosCsr, CSRSetting.CsrPath)
		if err != nil {
			panic(err)
		}
		err = tools.CreateFile(defaultIosKey, CSRSetting.KeyPath)
		if err != nil {
			panic(err)
		}
	}
}

func createPath(path string) {
	if !isExist(path) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			zap.S().Errorf("setting.Setup, fail to mkdir %s: %v", path, err)
		}
	}
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}

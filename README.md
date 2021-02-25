# README.md

## 这是什么
一个用go实现的iOS重签名模块，即市面上的iOS超级签名、蒲公英ios内测分发原理

使用本模块可以进行基本的IPA安装包重签名分发

实现功能：苹果开发者账号管理、IPA安装包管理

## 前提
1.生成ios.csr和ios.key文件
```bash
openssl genrsa -out ios.key 2048
openssl req -new -sha256 -key ios.key -out ios.csr
```
2.需要https，自行配置ssl证书

## 手动部署

```bash
git clone https://github.com/togettoyou/super-signature.git
# 进入项目isign目录下
cd super-signature/isign
# 更改Python默认编码
cp ./sitecustomize.py /usr/lib/python2.7/sitecustomize.py
# 验证(输出utf-8即更改成功)
python -c "import sys; print sys.getdefaultencoding()"
# 安装isign
tar xvf isign.tar.gz
cd isign
./version.sh
python setup.py build
python setup.py install
# 验证
isign -h
# 开启go mod
go env -w GO111MODULE=on
# 回到项目目录(记得更改app.ini配置信息)
cd super-signature
go run main.go
# 浏览器访问 http://localhost:10016/swagger/index.html
```

## 使用docker部署
```bash
docker-compose up
```

详见`Dockerfile`和`docker-compose.yml`文件

## 原理
[语雀浏览](https://www.yuque.com/togettoyou/cjqm/rbk50t)

### 基本流程
1. 添加Apple开发者账号(绑定App Store Connect API)
1. 根据描述文件获得用户设备的UDID
1. 借助App Store Connect API在开发者中心添加UDID、创建证书等
1. 重签名（使用isign开源项目实现在linux服务器上重签名） 
1. 将ipa包上传到服务器上，配置itms-service服务来做分发



> API

![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1614157707280-fc55e268-dc64-4a95-ade2-fb15da135562.png#align=left&display=inline&height=295&margin=%5Bobject%20Object%5D&name=image.png&originHeight=884&originWidth=3294&size=165697&status=done&style=none&width=1098)

`/api/v1/getAllPackage` 返回数据格式
```json
{
  "code": 0,
  "msg": "成功",
  "data": [
    {
      "ID": 1,
      "IconLink": "应用图标地址",
      "BundleIdentifier": "应用包名",
      "Name": "应用名称",
      "Version": "应用版本号",
      "BuildVersion": "应用BuildVersion",
      "MiniVersion": "最低支持ios版本",
      "Summary": "简介",
      "AppLink": "应用下载地址，iPhone使用Safari浏览器访问即可下载",
      "Size": "应用大小",
      "Count": "累计下载量"
    }
  ]
}
```
> 界面截图

![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1614159853374-673e82af-a2f2-479d-9ef8-03da193ed801.png#align=left&display=inline&height=478&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1970&originWidth=1154&size=504847&status=done&style=none&width=280)
![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1614159891066-ce68c4bf-60d7-49ae-864b-270e67a72a74.png#align=left&display=inline&height=477&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1964&originWidth=1154&size=691524&status=done&style=none&width=280)


### 添加Apple开发者账号

> API接口文档：[https://developer.apple.com/documentation/appstoreconnectapi](https://developer.apple.com/documentation/appstoreconnectapi)



使用App Store Connect API需要到[https://appstoreconnect.apple.com/access/api](https://appstoreconnect.apple.com/access/api) 生成API密钥P8文件，
以及对应的密钥ID和账号的Issuer ID。
![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1614157937920-e048fc1b-b8ef-4b08-a559-bcf0a9b72c39.png#align=left&display=inline&height=323&margin=%5Bobject%20Object%5D&name=image.png&originHeight=970&originWidth=3284&size=177328&status=done&style=none&width=1094.6666666666667)

首先创建一个结构体存放p8(下载的API密钥文件内容)，kid (密钥ID)，Iss (Issuer ID)
```go
type Authorize struct {
	P8  string
	Iss string
	Kid string
}
```
使用fasthttp来发送http请求
```go
func (a Authorize) httpRequest(method string, url string, body []byte) (*fasthttp.Response, error) {
	token, err := a.createToken()
	if err != nil {
		return nil, err
	}
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetContentType("application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	req.Header.SetMethod(method)
	req.SetRequestURI(url)
	req.SetBody(body)
	if err := fasthttp.Do(req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
```
token验证关键代码
```go
func (a Authorize) createToken() (string, error) {
	token := &jwt.Token{
		Header: map[string]interface{}{
			"alg": "ES256",
			"kid": a.Kid,
		},
		Claims: jwt.MapClaims{
			"iss": a.Iss,
			"exp": time.Now().Add(time.Second * 60 * 5).Unix(),
			"aud": "appstoreconnect-v1",
		},
		Method: jwt.SigningMethodES256,
	}
	privateKey, err := authKeyFromBytes([]byte(a.P8))
	if err != nil {
		return "", err
	}
	return token.SignedString(privateKey)
}

func authKeyFromBytes(key []byte) (*ecdsa.PrivateKey, error) {
	var err error
	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, errors.New("token: AuthKey must be a valid .p8 PEM file")
	}
	// Parse the key
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
		return nil, err
	}
	var pkey *ecdsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*ecdsa.PrivateKey); !ok {
		return nil, errors.New("token: AuthKey must be of type ecdsa.PrivateKey")
	}
	return pkey, nil
}

```
调用
```go
resp, err := authorize.httpRequest("GET", "https://api.appstoreconnect.apple.com/v1/devices", nil)
defer fasthttp.ReleaseResponse(resp)
if err != nil {
	return err
}
```
这样就可以直接借助App Store Connect API来完成添加udid、创建Certificates证书、创建BundleIds、创建Profile等来实现超级签名的核心功能。


### 获取UDID

#### 一、创建udid.mobileconfig文件
```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
    <dict>
        <key>PayloadContent</key>
        <dict>
            <key>URL</key>
            <string>https://qq.com</string> //回调接收UDID等信息的，借用这个回调地址将udid传到服务器后台
            <key>DeviceAttributes</key>
            <array>
                <string>UDID</string>
                <string>IMEI</string>
                <string>ICCID</string>
                <string>VERSION</string>
                <string>PRODUCT</string>
            </array>
        </dict>
        <key>PayloadOrganization</key>
        <string>仅用于查询设备UDID安装APP</string>
        <key>PayloadDisplayName</key>
        <string>仅用于查询设备UDID安装APP</string>
        <key>PayloadVersion</key>
        <integer>1</integer>
        <key>PayloadUUID</key>
        <string>62b2dc10-72be-4022-bdec-f1ea3c720d1a</string> //可在https://www.guidgen.com/随机生成
        <key>PayloadIdentifier</key>
        <string>com.yy.UDID-server</string>
        <key>PayloadDescription</key>
        <string>仅用于查询设备UDID安装APP</string>
        <key>PayloadType</key>
        <string>Profile Service</string>
    </dict>
</plist>
```
> iPhone使用Safari浏览器访问放在服务器上的mobileconfig文件，进行安装描述文件，安装完成后苹果会回调我们设置的url，就可以得到udid信息。设置的url是一个post接口，接收到udid信息处理完逻辑后，301重定向到我们需要跳转的网站，如果不301重定向，iPhone会显示安装失败！

![image.png](https://cdn.nlark.com/yuque/0/2020/png/1077776/1590130668052-c3862520-1e6c-4b1f-bf6e-2f57074ee225.png#align=left&display=inline&height=514&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1334&originWidth=750&size=122560&status=done&style=none&width=289)
#### 二、解析苹果返回的Plist信息，提取UDID
```xml
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>IMEI</key>
    <string>12 345678 901234 566789</string>
    <key>PRODUCT</key>
    <string>iPhone10,3</string>
    <key>UDID</key>
    <string>abcd0123456789XXXXXXXXXXXX</string>
    <key>VERSION</key>
    <string>12345</string>
  </dict>
</plist>
```
只需要解析出udid，调用App Store Connect API将UDID添加到苹果开发者中心即可。

### 重签名

> 添加开发者账号之前在本地使用openssl生成后续所需要的key和csr文件

```bash
openssl genrsa -out ios.key 2048
openssl req -new -sha256 -key ios.key -out ios.csr
```
> 利用csr文件调用CreateCertificates (App Store Connect API) 可以生成cer 证书
> 
> 接着利用cer证书生成pem文件（公钥）

```bash
openssl x509 -in ios_development.cer -inform DER -outform PEM -out ios_development.pem
```
> 公钥ios_development.pem、私钥ios.key、描述文件mobileprovision（调用CreateProfile App Store Connect API）、原始ipa
> 四大材料已凑齐！
> 

> 使用开源项目isign实现在Linux重签名(得到新的重签名安装包new.ipa)

```bash
isign -c ios_development.pem -k ios.key -p 描述文件.mobileprovision  -o new.ipa Runner.ipa
```


附：isign安装教程([https://github.com/apperian/isign](https://github.com/apperian/isign))
```bash
./version.sh
python setup.py build
python setup.py install
```
[isign.zip](https://www.yuque.com/attachments/yuque/0/2020/zip/1077776/1591342709219-ffa8ad4c-9f96-45c1-8b99-34965b58ec4c.zip?_lake_card=%7B%22uid%22%3A%221591342692957-0%22%2C%22src%22%3A%22https%3A%2F%2Fwww.yuque.com%2Fattachments%2Fyuque%2F0%2F2020%2Fzip%2F1077776%2F1591342709219-ffa8ad4c-9f96-45c1-8b99-34965b58ec4c.zip%22%2C%22name%22%3A%22isign.zip%22%2C%22size%22%3A34601410%2C%22type%22%3A%22application%2Fzip%22%2C%22ext%22%3A%22zip%22%2C%22progress%22%3A%7B%22percent%22%3A99%7D%2C%22status%22%3A%22done%22%2C%22percent%22%3A0%2C%22id%22%3A%22s3MfS%22%2C%22card%22%3A%22file%22%7D)


### IPA分发

创建后缀为plist的文件，内容如下：
```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
        <key>items</key>
        <array>
                <dict>
                        <key>assets</key>
                        <array>
                                <dict>
                                    <key>kind</key>
                                    <string>software-package</string>
                                    <key>url</key>
                                    <string>https://重签名后的ipa下载地址</string>
                                </dict>
                        </array>
                        <key>metadata</key>
                        <dict>
                            <key>bundle-identifier</key>
                            <string>com.togettoyou.app</string>
                            <key>bundle-version</key>
                            <string>1.0.0</string>
                            <key>kind</key>
                            <string>software</string>
                            <key>title</key>
                            <string>App</string>
                        </dict>
                </dict>
        </array>
</dict>
</plist>
```
安装用户需在Safari浏览器访问如下html：
```xml
<a href="itms-services://?action=download-manifest&url={{ .plist下载地址 }}">安装APP</a>
```
![image.png](https://cdn.nlark.com/yuque/0/2020/png/1077776/1590131104699-a9a81563-5b48-44b4-ab25-9a76353195ae.png#align=left&display=inline&height=496&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1334&originWidth=750&size=196955&status=done&style=none&width=279)


### 常见问题

1. 使用进行重签名时，如下报错：
```shell
Traceback (most recent call last):
  File "/bin/isign", line 4, in <module>
    __import__('pkg_resources').run_script('isign==1.6.15.1589436891.dev72+root', 'isign')
  File "/usr/lib/python2.7/site-packages/pkg_resources/__init__.py", line 666, in run_script
    self.require(requires)[0].run_script(script_name, ns)
  File "/usr/lib/python2.7/site-packages/pkg_resources/__init__.py", line 1462, in run_script
    exec(code, namespace, namespace)
  File "/usr/lib/python2.7/site-packages/isign-1.6.15.1589436891.dev72+root-py2.7.egg/EGG-INFO/scripts/isign", line 199, in <module>
    isign.resign(app_path, **kwargs)
  File "/usr/lib/python2.7/site-packages/isign-1.6.15.1589436891.dev72+root-py2.7.egg/isign/isign.py", line 79, in resign
    alternate_entitlements_path)
  File "/usr/lib/python2.7/site-packages/isign-1.6.15.1589436891.dev72+root-py2.7.egg/isign/archive.py", line 397, in resign
    ua = archive.unarchive_to_temp()
  File "/usr/lib/python2.7/site-packages/isign-1.6.15.1589436891.dev72+root-py2.7.egg/isign/archive.py", line 262, in unarchive_to_temp
    process_watchkit(app_dir, REMOVE_WATCHKIT)
  File "/usr/lib/python2.7/site-packages/isign-1.6.15.1589436891.dev72+root-py2.7.egg/isign/archive.py", line 78, in process_watchkit
    watchkit_paths = get_watchkit_paths(root_bundle_path)
  File "/usr/lib/python2.7/site-packages/isign-1.6.15.1589436891.dev72+root-py2.7.egg/isign/archive.py", line 63, in get_watchkit_paths
    bundle = Bundle(path)
  File "/usr/lib/python2.7/site-packages/isign-1.6.15.1589436891.dev72+root-py2.7.egg/isign/bundle.py", line 55, in __init__
    log.debug(u"Missing/invalid CFBundleSupportedPlatforms value in {}".format(self.info_path))
UnicodeDecodeError: 'ascii' codec can't decode byte 0xe9 in position 26: ordinal not in range(128)
```
这是Python的默认编码为ascii导致的，只需要改为utf-8即可解决。
直接在`/usr/lib/python2.7/site-packages`(具体路径结合实际情况)下新建一个`sitecustomize.py`文件(若不存在`site-packages`目录，则直接在`/usr/lib/python2.7`下创建)：
```python
# sitecustomize.py
# this file can be anywhere in your Python path,
# but it usually goes in ${pythondir}/lib/site-packages/
import sys
sys.setdefaultencoding('utf-8')
```
![image.png](https://cdn.nlark.com/yuque/0/2020/png/1077776/1590652483792-7a374715-9187-41b0-a1ea-f37b3881ebfd.png#align=left&display=inline&height=315&margin=%5Bobject%20Object%5D&name=image.png&originHeight=630&originWidth=612&size=39613&status=done&style=none&width=306)

验证是否成功：
命令行输入`python -c "import sys; print sys.getdefaultencoding()"`
![image.png](https://cdn.nlark.com/yuque/0/2020/png/1077776/1590652388536-09409bb1-5ecf-4d5a-ad21-1e7544bcb5fd.png#align=left&display=inline&height=23&margin=%5Bobject%20Object%5D&name=image.png&originHeight=46&originWidth=867&size=7550&status=done&style=none&width=433.5)

# README.md

![扫码_搜索联合传播样式-标准色版.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1612960247290-a878d022-cdd1-4f8b-ad39-98bafbe48894.png#align=left&display=inline&height=624&margin=%5Bobject%20Object%5D&name=%E6%89%AB%E7%A0%81_%E6%90%9C%E7%B4%A2%E8%81%94%E5%90%88%E4%BC%A0%E6%92%AD%E6%A0%B7%E5%BC%8F-%E6%A0%87%E5%87%86%E8%89%B2%E7%89%88.png&originHeight=624&originWidth=2092&size=5221770&status=done&style=none&width=2092#id=DA281&originHeight=624&originWidth=2092&originalType=binary&status=done&style=none#id=GqhuE&originHeight=624&originWidth=2092&originalType=binary&ratio=1&status=done&style=none)

# JetBrains 开源证书支持
感谢 JetBrains 提供的免费授权

<a href="https://www.jetbrains.com/?from=togettoyou" target="_blank"><img src="https://user-images.githubusercontent.com/55381228/127271051-14879011-41dd-4d1b-88a2-1591925b51de.png" width="250" align="middle"/></a>



## 这是什么

一个用go实现的iOS重签名应用，即市面上的iOS超级签名、蒲公英ios内测分发原理

使用本模块可以进行基本的IPA安装包重签名分发

实现功能：苹果开发者账号管理、IPA安装包管理

运行环境：Linux

## 代码架构使用
https://github.com/togettoyou/go-one-server

## 前提（重要，重要，重要）

1.生成替换 ios.csr 和 ios.key 文件

```bash
openssl genrsa -out ios.key 2048
openssl req -new -sha256 -key ios.key -out ios.csr
```

2.需要 https（获取UUID过程苹果服务器会回调我们的接口，需要https），可直接通过 config.yaml 开启 enableHttps （需要配置本地ssl证书），或通过 nginx 等网关代理等形式部署证书。

3.更改 config.yaml 配置 mysql.dsn（默认连接docker启动的mysql），applePath.url（必须改为你自己的https域名） 等信息

## 使用docker一键部署（推荐，推荐，推荐）

```bash
docker-compose up
# 会拉取 isign 签名环境，自动创建数据库，编译运行
```

详见`Dockerfile`和`docker-compose.yml`文件

部署截图：

仅作为例子，我的部署环境为本地虚机，https 域名使用内网穿透方式实现。

![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1622719557784-d384350b-d581-4784-b655-235f901ca571.png#clientId=uc4af6cdf-c3d2-4&from=paste&height=602&id=fnZsv&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1204&originWidth=1604&originalType=binary&size=128584&status=done&style=none&taskId=u01cf6afb-61d5-4cbe-b4f2-6a7670beac7&width=802#id=soLIX&originHeight=1204&originWidth=1604&originalType=binary&ratio=1&status=done&style=none)

更改 config.yaml 域名，启动服务：

![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1622719420286-2f0d32cb-0c48-48fb-a65a-c243643ef659.png#clientId=uc4af6cdf-c3d2-4&from=paste&height=716&id=ue315c056&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1432&originWidth=2353&originalType=binary&size=222303&status=done&style=none&taskId=u5e46eac2-7afe-4217-9224-45fcb15d221&width=1176.5#id=bhoQR&originHeight=1432&originWidth=2353&originalType=binary&ratio=1&status=done&style=none)
![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1622719722870-c1123242-92b7-4b03-a581-a06234b31892.png#clientId=uc4af6cdf-c3d2-4&from=paste&height=718&id=ucf800682&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1436&originWidth=2353&originalType=binary&size=395443&status=done&style=none&taskId=u2c2233e2-4124-47f2-80ac-9a4e5bafefd&width=1176.5#id=a53p6&originHeight=1436&originWidth=2353&originalType=binary&ratio=1&status=done&style=none)

访问 [https://isign.cn.utools.club/swagger/index.html](https://isign.cn.utools.club/swagger/index.html)

![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1622719814015-5552a7a4-496a-4271-b43f-7f78592176d1.png#clientId=uc4af6cdf-c3d2-4&from=paste&height=827&id=u84b71819&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1654&originWidth=2880&originalType=binary&size=275056&status=done&style=none&taskId=ua10a445e-d046-46a0-b6ef-617fde81539&width=1440#id=PdB8i&originHeight=1654&originWidth=2880&originalType=binary&ratio=1&status=done&style=none)

## 手动部署（需手动部署 isign 环境，略，自行百度谷歌 isign 部署）

```bash
# 部署isign环境略，验证安装是否成功
isign -h

git clone https://github.com/togettoyou/super-signature.git
cd super-signature
# 配置 config.yaml
go run main.go
# 浏览器访问 http://localhost:8888/swagger/index.html
```

## 使用说明

1、 上传苹果开发者账号信息

登陆 [https://appstoreconnect.apple.com/access/api](https://appstoreconnect.apple.com/access/api) 获取p8(下载的API密钥文件内容)，kid (
密钥ID)，Iss (Issuer ID)：

![](https://cdn.nlark.com/yuque/0/2021/png/1077776/1614157937920-e048fc1b-b8ef-4b08-a559-bcf0a9b72c39.png?x-oss-process=image%2Fwatermark%2Ctype_d3F5LW1pY3JvaGVp%2Csize_14%2Ctext_Z2l0aHViL3RvZ2V0dG95b3U%3D%2Ccolor_FFFFFF%2Cshadow_50%2Ct_80%2Cg_se%2Cx_10%2Cy_10#from=url&id=ipJUH&margin=%5Bobject%20Object%5D&originHeight=970&originWidth=3284&originalType=binary&ratio=2&status=done&style=none)

上传：

![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1623042480919-37ecee18-c7e7-4e17-91ac-c2ad8b7e117a.png#clientId=uab37fe2a-4554-4&from=paste&height=821&id=u8d372f30&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1641&originWidth=2880&originalType=binary&ratio=2&size=239573&status=done&style=none&taskId=ueb474557-a63b-43a0-97c2-b066502a2a4&width=1440)

2、 上传IPA

![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1623042643053-67a10d99-3359-4ebb-9ee4-b36d7ea48bdb.png#clientId=uab37fe2a-4554-4&from=paste&height=822&id=udac83704&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1644&originWidth=2880&originalType=binary&ratio=2&size=240127&status=done&style=none&taskId=ub147db0c-bab9-4419-abef-6de3e71fb46&width=1440)

3、 iPhone 使用 Safari 浏览器打开 AppLink 链接

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

![](https://cdn.nlark.com/yuque/0/2021/png/1077776/1614159853374-673e82af-a2f2-479d-9ef8-03da193ed801.png#from=url&id=yGJKs&margin=%5Bobject%20Object%5D&originHeight=1970&originWidth=1154&originalType=binary&ratio=2&status=done&style=none)

## 详细原理说明

[语雀浏览](https://www.yuque.com/togettoyou/cjqm/rbk50t)

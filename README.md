## 微信公众号：gopher云原生

go、云原生技术 ...... 来者不拒，欢迎关注

<img src="https://user-images.githubusercontent.com/55381228/221747734-13783ce6-1969-4c10-acd6-833f5046aa85.png" width="300px">

## 起源

每个客户端开发者都会想独立开发一款自己的 APP 。 但在开发过程中， iOS 不像 Android 那样可以自由的分享应用给其他人安装体验 （Android 只要把 `.apk` 发过去就行了）

对于 iOS 开发者来说， Apple 的个人开发者账号几乎是必备的， 而 Apple 公司允许我们添加 100 台设备绑定到账号上，这 100 台设备可以自由安装由该账号签名且使用 Ad Hoc 方式打包出的 `.ipa`

本项目就是利用这个官方的规则，来简化 iOS APP 的分享流程。当你开发一款 APP 的过程中，想要给身边的小伙伴体验一下，只需要部署本项目，然后上传你的个人开发者账号和 `.ipa` 就可以生成一个二维码，扫一扫即可全程自动绑定设备并签名安装 APP

## 原理

设备绑定账号调用 [App Store Connect API](https://developer.apple.com/documentation/appstoreconnectapi) 实现

`.ipa` 签名调用 [zsign](https://github.com/zhlynn/zsign) 实现

## 注意事项

本项目完全开源免费，纯技术分享，不提供任何平台支持，仅作为给个人开发者分享自己的测试 APP 使用，严禁使用本项目进行任何商业盈利、损害官方利益、分享任何违法违规的 APP 等行为

## 个人开发者账号

本项目添加开发者账号后会占用账号的一个 iOS Development certificate 名额（每个账号最多只能创建两个），所以你可能得预留一个

其中添加开发者账号所需的 iss（Issuer ID）、kid（密钥 ID）、P8 文件（API 密钥）需要在 https://appstoreconnect.apple.com/access/api 创建后获取

## Railway 部署

使用 Railway 一键部署，每月有 5 美元的免费额度：

[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/new/template?code=xOgXXB&referralCode=FVN0mI)

## Docker 部署

创建 `docker-compose.yml` ，内容如下：

```yaml
version: '3'
services:
  server:
    image: togettoyou/ipashare:latest
    ports:
      - "8888:8888"
    volumes:
      - $PWD/data:/root/togettoyou/data
      - /etc/timezone:/etc/timezone
      - /etc/localtime:/etc/localtime
    environment:
      - SERVER_URL=https://example.com # 必须修改，域名
      - SERVER_HTTPPORT=8888 # 指定后端服务端口
      - SERVER_MAXJOB=10 # 签名任务并发量
      - SERVER_RUNMODE=release
      - LOG_LEVEL=info
```

执行：

```shell
# 前台运行服务
docker-compose up
# 后台运行服务
docker-compose up -d
# 查看服务日志
docker-compose logs -f --tail 100
# 停止服务
docker-compose stop
# 启动服务
docker-compose start
# 重启服务
docker-compose restart
# 删除服务
docker-compose down -v
```

启动成功后可通过 `域名/admin` 例如 `https://example.com/admin` 访问后台管理

项目启动后所有数据保存在 `data` 目录，请妥善保管

- `apple_developer` : 苹果开发者账号相关证书
- `temporary_file_path` : 临时文件存放路径（打包后的 IPA，1小时后会自动删除）
- `upload_file_path` : 文件上传路径（上传的 IPA）
- `sqlite.db` : 默认的sqlite数据库文件

数据库支持更换为 MySQL ，在 `docker-compose.yml` 加入以下环境变量即可：

```yaml
- MYSQL_ENABLE=true
- MYSQL_DSN=root:123456@tcp(127.0.0.1:3306)/db_default?charset=utf8mb4&parseTime=True&loc=Local
```

更多环境变量可参考：[server/conf/default.yaml](server/conf/default.yaml) ，变量层级使用 `_` 连接，如 `MYSQL_MAXIDLE` 代表 mysql
配置中的空闲连接池中连接的最大数量

## JetBrains 开源证书支持

本项目使用 GoLand 开发，感谢 JetBrains 提供的免费授权

<a href="https://www.jetbrains.com/?from=togettoyou" target="_blank"><img src="https://user-images.githubusercontent.com/55381228/127271051-14879011-41dd-4d1b-88a2-1591925b51de.png" width="250" align="middle"/></a>

## 部分效果预览

| ![login](https://user-images.githubusercontent.com/55381228/195557740-3b65e5c9-b86e-42ba-929e-273b0e110d23.png) | ![developer](https://user-images.githubusercontent.com/55381228/195557833-ec3d4db8-76ee-4d60-9915-ee35f06f2efe.png) |
| :----------------------------------------------------------: | :----------------------------------------------------------: |
|                         后台管理登录                         |                        开发者账号管理                        |
| ![ipalist](https://user-images.githubusercontent.com/55381228/195557932-54b8ca9b-081d-4ddf-bbd7-5b6004664720.png) | ![keylist](https://user-images.githubusercontent.com/55381228/195558156-7b7dea93-d9d6-4aac-b0a9-3bf9751828d2.png) |
|                           应用管理                           |                         密钥管理                         |
| ![oss](https://user-images.githubusercontent.com/55381228/232664237-13e74612-23cd-4dea-aac9-8b9b92c2cf2e.png) | ![mobileconfig](https://user-images.githubusercontent.com/55381228/232421233-b41de68d-5d78-4412-a1f7-cf69db356cdf.png) |
|                       阿里云 OSS 设置                        |                        描述文件签名                        |
| ![pw](https://user-images.githubusercontent.com/55381228/232421631-16c7b41d-1749-4c0f-b096-b894bf750416.png) | ![image](https://user-images.githubusercontent.com/55381228/232664767-3d50b491-e25d-46d6-8d19-6c5d302d7bab.png) |
|                       更改用户名密码                        |                           安装 APP                     |


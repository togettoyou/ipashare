## 微信公众号：SuperGopher

go、云原生技术、项目问题、单纯支持 ...... 来者不拒
<div align="center" style="width: 50%">

![微信公众号.png](https://user-images.githubusercontent.com/55381228/155444889-eacc0104-cd85-45c9-b7b7-9036e0c2334c.jpg)
</div>

## 起源

每个客户端开发者都会想独立开发一款自己的 APP 。 但 iOS 不像 Android 那样可以自由分发应用 （Android 只要把 apk 甩出去就行了）

对于 iOS 开发者来说，苹果开发者账号几乎人手必备（这里只讨论个人账号）， 而苹果公司允许我们添加 100 台设备（udid）绑定到账号上，这 100 台设备可以自由安装由账号签名且使用 Ad Hoc 方式打包出的 `.ipa`

本项目就是利用这个规则，来简化 iOS APP 的分发流程。当你开发一款 APP 的过程中，想要给身边的小伙伴体验一下，只需要使用本项目生成一个二维码链接，扫一扫，即可全程自动绑定设备并签名安装 APP

## 注意事项

本项目核心功能调用 [zsign](https://github.com/zhlynn/zsign)
和 [App Store Connect API](https://developer.apple.com/documentation/appstoreconnectapi) 实现

本项目添加开发者账号后会占用账号的一个 iOS Development certificate 名额（每个账号最多只能创建两个），所以你可能得预留一个

本项目仅作为给开发者分发合法合规的 APP 使用，严禁使用本项目进行任何盈利、损害官方利益、分发任何违法违规的 APP 等行为

## 效果预览

[https://supersign.up.railway.app/admin](https://supersign.up.railway.app/admin)

开发者账号管理

![image](https://user-images.githubusercontent.com/55381228/157637112-9ea00b7b-9edb-441e-bed9-c70a8252db62.png)

开发者账号的 iss、kid、P8 文件需要到 https://appstoreconnect.apple.com/access/api 生成

![image](https://user-images.githubusercontent.com/55381228/157636136-b0705dfe-5887-40a5-a464-f736157f5335.png)

应用管理

![image](https://user-images.githubusercontent.com/55381228/157637583-6ecb301f-490d-48cb-93f7-eb53c6c7ba9a.png)

支持阿里云OSS设置

![image](https://user-images.githubusercontent.com/55381228/157637721-535976b5-9110-4334-905c-cab03fcd0402.png)

用户名密码修改

![image](https://user-images.githubusercontent.com/55381228/157637800-b40e80cd-eaf6-4603-a212-66fde299d07e.png)

## 部署项目

提供了多种部署方案，以供参考，不同方案的区别仅在于 `docker-compose.yml` 的配置不同，部署时根据实际情况选择其中任意一种即可

### 配置 1

需要准备的东西：

1. 一个域名，并解析到你的服务器IP地址上，例如：`example.com`
2. 服务器开放 `80` 和 `443` 端口，并确保没有被其它程序占用
3. 一个邮箱地址，用于自动向 Let's Encrypt 申请免费SSL证书，例如：`foo@qq.com`

创建 `docker-compose.yml` ，内容如下：

```yaml
version: '3'
services:
  server:
    image: togettoyou/supersign:latest
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - $PWD/data:/root/togettoyou/data
      - /etc/timezone:/etc/timezone
      - /etc/localtime:/etc/localtime
    environment:
      - SERVER_URL=https://example.com
      - SERVER_TLS=true
      - SERVER_AUTOTLS=true
      - SERVER_ACMEEMAIL=foo@qq.com
      - SERVER_MAXJOB=10
      - SERVER_RUNMODE=release
      - LOG_LEVEL=info
```

优点：快速部署，自动配置SSL证书

缺点：需要占用 `80` 和 `443` 端口

### 配置 2

需要准备的东西：

1. 一个域名，并解析到你的服务器IP地址上，例如：`example.com`
2. 服务器开放 `80` 和 `443` 端口，并确保没有被其它程序占用
3. 自行准备域名的SSL证书，例如：`ssl.crt` 和 `ssl.key`

创建 `docker-compose.yml` ，内容如下：

```yaml
version: '3'
services:
  server:
    image: togettoyou/supersign:latest
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - $PWD/data:/root/togettoyou/data
      - $PWD/ssl.crt:/root/togettoyou/ssl.crt
      - $PWD/ssl.key:/root/togettoyou/ssl.key
      - /etc/timezone:/etc/timezone
      - /etc/localtime:/etc/localtime
    environment:
      - SERVER_URL=https://example.com
      - SERVER_TLS=true
      - SERVER_AUTOTLS=false
      - SERVER_CRT=ssl.crt
      - SERVER_KEY=ssl.key
      - SERVER_MAXJOB=10
      - SERVER_RUNMODE=release
      - LOG_LEVEL=info
```

优点：快速部署

缺点：需要占用 `80` 和 `443` 端口，需要自己准备SSL证书

### 配置 3

需要准备的东西：

1. 一个域名，并解析到你的服务器IP地址上，例如：`example.com`
2. 指定一个任意空闲端口，例如：`8888`
3. 服务器需要部署其它能够支持反向代理和SSL的网关服务，例如：Nginx

创建 `docker-compose.yml` ，内容如下：

```yaml
version: '3'
services:
  server:
    image: togettoyou/supersign:latest
    ports:
      - "8888:8888"
    volumes:
      - $PWD/data:/root/togettoyou/data
      - /etc/timezone:/etc/timezone
      - /etc/localtime:/etc/localtime
    environment:
      - SERVER_URL=https://example.com
      - SERVER_HTTPPORT=8888
      - SERVER_MAXJOB=10
      - SERVER_RUNMODE=release
      - LOG_LEVEL=info
```

此方案需要你额外在网关服务上配置反向代理以及SSL证书，以 Nginx 配置为例：

```nginx configuration
server {
    listen 80;
    server_name example.com;
    rewrite ^(.*)$ https://$host$1 permanent;
}

server {
    listen 443 ssl;
    server_name example.com;
    client_max_body_size 300m;
    client_body_buffer_size 50m;
    location / {
        proxy_pass http://127.0.0.1:8888;
        proxy_http_version 1.1;
        proxy_set_header Host $proxy_host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    ssl on;
    ssl_certificate /etc/nginx/ssl/domain.crt;
    ssl_certificate_key /etc/nginx/ssl/domain.key;
    ssl_session_timeout 5m;
    ssl_protocols TLSv1 TLSv1.1 TLSv1.2;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE:ECDH:AES:HIGH:!NULL:!aNULL:!MD5:!ADH:!RC4:!DH:!DHE;
    ssl_prefer_server_ciphers on;
}
```

优点：不需要占用 `80` 和 `443` 端口、可利用网关强大的能力（上传大小限制等）

### 启动和停止

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

### 持久化

项目启动后所有文件保存在 `data` 目录，请妥善保管

- `apple_developer` : 苹果开发者账号相关证书
- `temporary_file_path` : 临时文件存放路径（重签名后的 IPA）
- `upload_file_path` : 文件上传路径（上传的 IPA）
- `sqlite.db` : 默认的sqlite数据库文件

数据库支持更换为 MySQL ，在 `docker-compose.yml` 加入以下环境变量 ：

```yaml
- MYSQL_ENABLE=true
- MYSQL_DSN=root:123456@tcp(127.0.0.1:3306)/db_default?charset=utf8mb4&parseTime=True&loc=Local
```

更多环境变量可参考：[server/conf/default.yaml](server/conf/default.yaml) ，变量层级使用 `_` 连接，如 `MYSQL_MAXIDLE` 代表 mysql
配置中的空闲连接池中连接的最大数量

## 喝杯奶茶

如有帮助，可以打赏支持，一分也是爱！

|  ![微信打赏](https://user-images.githubusercontent.com/55381228/155450359-0ce92911-fd3f-4d6b-9878-e40a17b34652.jpg)   | ![支付宝打赏](https://user-images.githubusercontent.com/55381228/155450383-509d0475-5497-4983-8583-137946b4d78e.jpg)  |
|  ----  | ----  |
| 微信  | 支付宝 |

## JetBrains 开源证书支持

本项目使用 GoLand 开发，感谢 JetBrains 提供的免费授权

<a href="https://www.jetbrains.com/?from=togettoyou" target="_blank"><img src="https://user-images.githubusercontent.com/55381228/127271051-14879011-41dd-4d1b-88a2-1591925b51de.png" width="250" align="middle"/></a>

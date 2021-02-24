# README.md

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

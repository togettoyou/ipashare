.PHONY: all docs linux linux-docs docker docker-docs run gotool clean help

# 生成的二进制文件名
BINARY_NAME="supersign"
# 生成的镜像名:版本
IMAGE_NAME="supersign:v1"
# 项目包名
MOD_NAME="supersign"
# run 参数
TARGET=$(out)

# 编译添加版本信息
versionDir = "${MOD_NAME}/pkg/version"
gitTag = $(shell if [ "`git describe --tags --abbrev=0 2>/dev/null`" != "" ];then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
buildDate = $(shell TZ=Asia/Shanghai date +%FT%T%z)
gitCommit = $(shell git log --pretty=format:'%H' -n 1)
gitTreeState = $(shell if git status|grep -q 'clean';then echo clean; else echo dirty; fi)
ldflags="-w -X ${versionDir}.gitTag=${gitTag} -X ${versionDir}.buildDate=${buildDate} -X ${versionDir}.gitCommit=${gitCommit} -X ${versionDir}.gitTreeState=${gitTreeState}"

# 执行make命令时所执行的所有命令
all: clean
	CGO_ENABLED=0 go build -v -ldflags ${ldflags} -o ${BINARY_NAME} cmd/server/main.go

docs: clean
	CGO_ENABLED=0 go build -tags "docs" -v -ldflags ${ldflags} -o ${BINARY_NAME} cmd/server/main.go

# 交叉编译linux amd64版本
linux: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags ${ldflags} -o ${BINARY_NAME} cmd/server/main.go

linux-docs: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags "docs" -v -ldflags ${ldflags} -o ${BINARY_NAME} cmd/server/main.go

# 编译生成docker镜像
docker:
	docker build -t ${IMAGE_NAME} -f Dockerfile .

docker-docs:
	docker build -t ${IMAGE_NAME} -f Dockerfile.Docs .

# 运行项目
run:
	go run -tags "docs" cmd/server/main.go $(TARGET)

# gotool工具
gotool:
    # 整理代码格式 && 代码静态检查
	cd cmd/server && gofmt -w . && go vet . | grep -v vendor;true

# 清理二进制文件
clean:
	@if [ -f ${BINARY_NAME} ] ; then rm ${BINARY_NAME} ; fi

# 帮助
help:
	@echo "make - 编译生成当前平台可运行的二进制文件(不带swagger文档)"
	@echo "make docs - 编译生成当前平台可运行的二进制文件(带swagger文档)"
	@echo "make linux - 交叉编译生成linux amd64可运行的二进制文件(不带swagger文档)"
	@echo "make linux-docs - 交叉编译生成linux amd64可运行的二进制文件(带swagger文档)"
	@echo "make docker - 编译生成docker镜像(不带swagger文档)"
	@echo "make docker-docs - 编译生成docker镜像(带swagger文档)"
	@echo "make run - 直接运行 Go 代码(带swagger文档)\nmake run out='-c config.yaml' - 指定配置文件"
	@echo "make gotool - 运行 Go 工具 'fmt' and 'vet'"
	@echo "make clean - 清理编译生成的二进制文件"

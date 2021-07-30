# 定义伪目标。不创建目标文件，而是去执行这个目标下面的命令。
.PHONY: all linux run gotool clean help

# 生成的二进制文件名
BINARY_NAME="super-signature-app"
MODULE_NAME="super-signature"
TARGET=$(out)

# 编译添加版本信息
versionDir = "${MODULE_NAME}/util/version"
gitTag = $(shell if [ "`git describe --tags --abbrev=0 2>/dev/null`" != "" ];then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
buildDate = $(shell TZ=Asia/Shanghai date +%FT%T%z)
gitCommit = $(shell git log --pretty=format:'%H' -n 1)
gitTreeState = $(shell if git status|grep -q 'clean';then echo clean; else echo dirty; fi)
ldflags="-w -X ${versionDir}.gitTag=${gitTag} -X ${versionDir}.buildDate=${buildDate} -X ${versionDir}.gitCommit=${gitCommit} -X ${versionDir}.gitTreeState=${gitTreeState}"

# 执行make命令时所执行的所有命令
all: clean
	go build -v -ldflags ${ldflags} -o ${BINARY_NAME} .

# 交叉编译linux amd64版本
linux: clean
	GOOS=linux GOARCH=amd64 go build -v -ldflags ${ldflags} -o ${BINARY_NAME} .

# gotool工具
gotool:
    # 整理代码格式
	gofmt -w .
    # 代码静态检查
	go vet . | grep -v vendor;true

# 清理二进制文件
clean:
	@if [ -f ${BINARY_NAME} ] ; then rm ${BINARY_NAME} ; fi

# 帮助
help:
	@echo "make - 编译生成当前平台可运行的二进制文件"
	@echo "make linux - 交叉编译生成linux amd64可运行的二进制文件"
	@echo "make gotool - 运行 Go 工具 'fmt' and 'vet'"
	@echo "make clean - 清理编译生成的二进制文件"

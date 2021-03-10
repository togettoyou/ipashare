FROM golang:1.14

RUN apt update \
    && apt-get -y install zip python-setuptools \
    && curl https://bootstrap.pypa.io/pip/2.7/get-pip.py -o get-pip.py \
    && python get-pip.py \
    && pip install cryptography

COPY . /root/togettoyou/super-signature
COPY ./isign/sitecustomize.py /usr/lib/python2.7/sitecustomize.py

WORKDIR /root/togettoyou/super-signature/isign
RUN tar xvf isign.tar.gz \
    && chmod -R 755 isign
WORKDIR /root/togettoyou/super-signature/isign/isign
RUN ./version.sh \
    && python setup.py build \
    && python setup.py install

ENV GOPROXY https://goproxy.cn,direct
WORKDIR /root/togettoyou/super-signature
RUN go env -w GO111MODULE=on \
    && go build -o "super-signature" .

EXPOSE 10016
ENTRYPOINT ["./super-signature"]

# 《====可在app.ini配置数据库服务器====》
# 《====直接运行以下命令====》
# 编译
# docker build -t super-signature:v1 .
# 查看生成镜像
# docker images
# 启动容器
# docker run -it -d --name super-signature -p 10016:10016 super-signature:v1
# 可进入容器 docker exec -it super-signature bash
# 查看日志 docker logs -f 容器ID
# 验证Python编码 python -c "import sys; print sys.getdefaultencoding()"
# 验证服务是否启动 ps -A | grep super-signature
# 浏览器访问 http://localhost:9999/swagger/index.html
FROM golang:1.16 AS builder
ENV GO111MODULE=on
ENV GOPROXY https://goproxy.cn,direct
COPY . /root/togettoyou/
WORKDIR /root/togettoyou/
RUN make

FROM togettoyou/zsign:latest AS zsign

FROM centos:7
COPY --from=builder /root/togettoyou/supersign /root/togettoyou/
COPY --from=builder /root/togettoyou/conf/ /root/togettoyou/conf/
COPY --from=zsign /zsign/zsign /bin/zsign
RUN yum install -y openssl openssl-devel unzip zip sqlite3 libsqlite3-dev
WORKDIR /root/togettoyou/
ENTRYPOINT ["./supersign"]
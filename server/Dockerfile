FROM golang:1.16 AS builder
RUN apt-get update && apt-get install -y \
    git make gcc \
    && rm -rf /var/lib/apt/lists/*
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct
COPY . /root/togettoyou/
WORKDIR /root/togettoyou/
RUN make

FROM togettoyou/zsign:latest AS zsign

FROM centos:7
COPY --from=builder /root/togettoyou/ipashare /root/togettoyou/
COPY --from=builder /root/togettoyou/conf/ /root/togettoyou/conf/
COPY --from=zsign /zsign/zsign /bin/zsign
RUN yum install -y openssl openssl-devel unzip zip
WORKDIR /root/togettoyou/
ENTRYPOINT ["./ipashare"]

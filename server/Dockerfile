FROM centos:7 AS builder
RUN yum install -y wget git make gcc \
    && wget https://studygolang.com/dl/golang/go1.16.6.linux-amd64.tar.gz \
    && tar -zxvf go1.16.6.linux-amd64.tar.gz -C /usr/local/
ENV GOROOT=/usr/local/go
ENV PATH=$PATH:$GOROOT/bin
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct
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
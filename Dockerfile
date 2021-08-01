FROM centos:7 AS builder
RUN yum install -y wget git make gcc \
    && wget https://studygolang.com/dl/golang/go1.16.6.linux-amd64.tar.gz \
    && tar -zxvf go1.16.6.linux-amd64.tar.gz -C /usr/local/
ENV GOROOT=/usr/local/go
ENV PATH=$PATH:$GOROOT/bin
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct
COPY . /root/super-signature
WORKDIR /root/super-signature
RUN make

FROM centos:7
COPY --from=builder /root/super-signature/super-signature-app /root/super-signature/
COPY zsign/zsign /usr/local/bin/
RUN yum install -y openssl openssl-devel unzip \
    && chmod +x /usr/local/bin/zsign
WORKDIR /root/super-signature
ENTRYPOINT ["./super-signature-app"]
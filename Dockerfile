FROM centos:7 AS builder-server
RUN yum install -y wget git make gcc \
    && wget https://studygolang.com/dl/golang/go1.16.6.linux-amd64.tar.gz \
    && tar -zxvf go1.16.6.linux-amd64.tar.gz -C /usr/local/
ENV GOROOT=/usr/local/go
ENV PATH=$PATH:$GOROOT/bin
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct
WORKDIR /root/togettoyou/
COPY server/. .
RUN make

FROM node:16.17.1-alpine AS builder-web
WORKDIR /app
COPY web/package*.json ./
RUN yarn install
COPY web/. .
RUN yarn run build:prod

FROM togettoyou/zsign:latest AS zsign

FROM centos:7
WORKDIR /root/togettoyou/
COPY --from=builder-server /root/togettoyou/ipashare ./
COPY --from=builder-server /root/togettoyou/conf/ ./conf/
COPY --from=builder-web /app/dist/ ./dist/
COPY --from=zsign /zsign/zsign /bin/zsign
RUN yum install -y openssl openssl-devel unzip zip
ENTRYPOINT ["./ipashare"]

FROM golang:1.14 AS builder
ENV GO111MODULE=on
ENV GOPROXY https://goproxy.cn,direct
COPY . /root/togettoyou/super-signature
WORKDIR /root/togettoyou/super-signature
RUN make docs

FROM togettoyou/isign:latest
COPY --from=builder /root/togettoyou/super-signature/super-signature-app /root/togettoyou/super-signature/
COPY --from=builder /root/togettoyou/super-signature/conf/ /root/togettoyou/super-signature/conf/
WORKDIR /root/togettoyou/super-signature
EXPOSE 8888
ENTRYPOINT ["./super-signature-app"]
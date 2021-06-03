FROM golang:1.14 AS builder
COPY . /root/togettoyou/super-signature
ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.io"
WORKDIR /root/togettoyou/super-signature
RUN make docs

FROM togettoyou/isign:latest
COPY --from=builder /root/togettoyou/super-signature/app /root/togettoyou/super-signature/
WORKDIR /root/togettoyou/super-signature
EXPOSE 8888
ENTRYPOINT ["./app"]
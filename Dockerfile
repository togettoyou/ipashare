FROM centos:7 AS builder
RUN rpm --import https://mirror.go-repo.io/centos/RPM-GPG-KEY-GO-REPO \
    && curl -s https://mirror.go-repo.io/centos/go-repo.repo | tee /etc/yum.repos.d/go-repo.repo \
    && yum install -y git make golang \
    && go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct
COPY . /root/super-signature
WORKDIR /root/super-signature
RUN make

FROM centos:7
COPY --from=builder /root/super-signature/super-signature-app /root/super-signature/
COPY router/templates/ /root/super-signature/router/templates/
COPY conf/ /root/super-signature/conf/
COPY zsign/zsign /usr/local/bin/
RUN yum install -y openssl openssl-devel \
    && chmod +x /usr/local/bin/zsign
WORKDIR /root/super-signature
ENTRYPOINT ["./super-signature-app"]
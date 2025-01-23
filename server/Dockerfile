FROM golang:1.16 AS builder
COPY . /root/togettoyou/
WORKDIR /root/togettoyou/
RUN make

FROM togettoyou/zsign:latest AS zsign

FROM centos:7
COPY --from=builder /root/togettoyou/ipashare /root/togettoyou/
COPY --from=builder /root/togettoyou/conf/ /root/togettoyou/conf/
COPY --from=zsign /zsign/zsign /bin/zsign
RUN sed -i 's|mirrorlist=http://mirrorlist.centos.org|#mirrorlist=http://mirrorlist.centos.org|' /etc/yum.repos.d/CentOS-Base.repo && \
    sed -i 's|#baseurl=http://mirror.centos.org|baseurl=http://mirrors.aliyun.com|' /etc/yum.repos.d/CentOS-Base.repo && \
    yum clean all && \
    yum makecache
RUN yum install -y openssl openssl-devel unzip zip
WORKDIR /root/togettoyou/
ENTRYPOINT ["./ipashare"]

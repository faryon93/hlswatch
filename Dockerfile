FROM golang:alpine as builder
MAINTAINER Maximilian Pachl <m@ximilian.info>
# setup the environment
ENV TZ=Europe/Berlin

# install dependencies
RUN apk --update --no-cache add git gcc musl-dev tzdata
WORKDIR /go/src/github.com/faryon93/hlswatch
ADD ./ ./

# build the go binary
RUN go get github.com/faryon93/hlswatch && \
    go build -v -o /tmp/hlswatch .

FROM alpine:latest
MAINTAINER Maximilian Pachl <m@ximilian.info>

# configuration and versions
ENV NGINX_VERSION="1.13.11"
ENV BUILD_TOOLS="g++ make pcre-dev openssl-dev unzip"
ENV RUNTIME_LIBS="openssl pcre"

# setup users
RUN adduser -D -u 1000 -g 'www' www

# setup the rootfs
ADD hlswatch.conf /etc/
ADD entry.sh /
COPY --from=builder /tmp/hlswatch /usr/sbin/hlswatch
RUN chmod 755 /entry.sh && \
    chmod 755 /usr/sbin/hlswatch

# start command
CMD ["/entry.sh"]

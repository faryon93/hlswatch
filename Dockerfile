FROM alpine:3.5
MAINTAINER Maximilian Pachl <m@ximilian.info>

# configuration and versions
ENV NGINX_VERSION="1.11.9"
ENV BUILD_TOOLS="g++ make pcre-dev openssl-dev unzip"
ENV RUNTIME_LIBS="openssl pcre"

# download the sources
ADD http://nginx.org/download/nginx-$NGINX_VERSION.tar.gz /tmp
ADD https://github.com/arut/nginx-rtmp-module/archive/master.zip /tmp/nginx-rtmp-master.zip

# compile and install nginx
RUN apk add --update $BUILD_TOOLS $RUNTIME_LIBS && \
	cd /tmp && \
	tar xzvf nginx-$NGINX_VERSION.tar.gz && \
	unzip nginx-rtmp-master.zip && \
	cd /tmp/nginx-$NGINX_VERSION && \
	./configure --prefix=/usr \
				--modules-path=/var/lib/nginx/modules \
				--conf-path=/etc/nginx/nginx.conf \
				--pid-path=/var/run/nginx.pid \
				--lock-path=/var/run/nginx.lock \
				--sbin-path=/usr/sbin/nginx \
				--error-log-path=/var/log/nginx/error.log \
				--http-log-path=/var/log/nginx/access.log \
				--http-client-body-temp-path=/tmp/client_body_temp \
				--http-proxy-temp-path=/tmp/proxy_temp \
				--user=www \
				--group=www \
				--add-module=/tmp/nginx-rtmp-module-master \
				--without-http_fastcgi_module \
				--without-http_uwsgi_module \
				--without-http_scgi_module \
				--with-http_ssl_module  \ 
				--with-http_v2_module && \
	make -j5 && \
	make install && \
	rm -r /usr/html && \

# remove build tools
	rm -r /tmp/nginx-$NGINX_VERSION && \
	rm -r /tmp/nginx-rtmp-module-master && \
	rm /tmp/nginx-rtmp-master.zip && \
	rm /tmp/nginx-$NGINX_VERSION.tar.gz && \
	apk del $BUILD_TOOLS && \
	rm -rf /var/cache/apk/*

# setup users
RUN adduser -D -u 1000 -g 'www' www

# network configuration
EXPOSE 1935
EXPOSE 80
EXPOSE 443

# setup the rootfs
ADD rootfs /
RUN mkdir /tmp/hls && \
    chmod 755 /usr/sbin/startup && \
    chmod 755 /usr/sbin/hlswatch

# start command
CMD ["/usr/sbin/startup"]

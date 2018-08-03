FROM alpine:3.8

ARG MYSQL_VERSION=${MYSQL_VERSION:-"10.2.15-r0"}
ARG MYSQL_CLIENT_VERSION=${MYSQL_CLIENT_VERSION:-"10.2.15-r0"}

RUN apk update && \
    apk add --no-cache --no-progress mysql=${MYSQL_VERSION} mysql-client=${MYSQL_CLIENT_VERSION} && \
    addgroup mysql mysql && \
    mkdir -p /opt/scripts && \
    rm -rf /var/cache/apk/*

VOLUME ["/var/lib/mysql", "/var/lib/mysql"]

COPY ./scripts/docker-entrypoint.sh /opt/scripts/docker-entrypoint.sh
RUN chmod +x /opt/scripts/docker-entrypoint.sh

EXPOSE 3306

ENTRYPOINT ["/opt/scripts/docker-entrypoint.sh"]
# CMD []
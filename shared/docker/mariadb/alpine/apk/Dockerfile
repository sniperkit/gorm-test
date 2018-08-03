# obtain latest alpine linux image
FROM alpine:3.8

ARG MARIADB_VERSION=${MARIADB_VERSION:-"10.2.15-r0"}
ARG MARIADB_CLIENT_VERSION=${MARIADB_CLIENT_VERSION:-"10.2.15-r0"}

# upgrade
# RUN apk add --update --no-cache --no-progress mariadb=${MARIADB_VERSION} mariadb-client=${MARIADB_CLIENT_VERSION} \
#	&& rm -fr /var/lib/apk/* \
#	&& rm -fr /var/cache/apk/* \
#	&& sed -Ei 's/^(bind-address|log)/#&/' /etc/mysql/my.cnf \
#	&& echo "skip-name-resolve" | awk '{ print } $1 == "[mysqld]" && c == 0 { c = 1; system("cat") }' /etc/mysql/my.cnf > /tmp/my.cnf && mv /tmp/my.cnf /etc/mysql/my.cnf \
#	&& echo "skip-host-cache" | awk '{ print } $1 == "[mysqld]" && c == 0 { c = 1; system("cat") }' /etc/mysql/my.cnf > /tmp/my.cnf && mv /tmp/my.cnf /etc/mysql/my.cnf \
#	&& echo "bind-address = 0.0.0.0" | awk '{ print } $1 == "[mysqld]" && c == 0 { c = 1; system("cat") }' /etc/mysql/my.cnf > /tmp/my.cnf && mv /tmp/my.cnf /etc/mysql/my.cnf \
#	&& echo "mysql_install_db --user=mysql" > /tmp/config \
# 	&& echo "mysqld_safe &" >> /tmp/config \
#  	&& echo "mysqladmin --silent --wait=30 ping || exit 1" >> /tmp/config \
#  	&& echo "mysql -uroot --execute=\"CREATE USER 'admin'@'%' IDENTIFIED BY 'admin';\"" >> /tmp/config \
#  	&& echo "mysql -uroot --execute=\"GRANT ALL PRIVILEGES ON *.* TO 'admin'@'%' WITH GRANT OPTION;\"" >> /tmp/config \
#  	&& echo "mysqladmin reload" >> /tmp/config \
#  	&& sh /tmp/config \
#  	&& rm -f /tmp/config

ENV LC_ALL=en_US.UTF-8

RUN mkdir /docker-entrypoint-initdb.d && \
    apk -U upgrade && \
    apk add --no-cache --no-cache --no-progress mariadb=${MARIADB_VERSION} mariadb-client=${MARIADB_CLIENT_VERSION} && \
    apk add --no-cache tzdata && \
    # clean up
    rm -rf /var/cache/apk/*

# comment out a few problematic configuration values
RUN sed -Ei 's/^(bind-address|log)/#&/' /etc/mysql/my.cnf && \
    # don't reverse lookup hostnames, they are usually another container
    sed -i '/^\[mysqld]$/a skip-host-cache\nskip-name-resolve' /etc/mysql/my.cnf && \
    # always run as user mysql
    sed -i '/^\[mysqld]$/a user=mysql' /etc/mysql/my.cnf && \
    # allow custom configurations
    echo -e '\n!includedir /etc/mysql/conf.d/' >> /etc/mysql/my.cnf && \
    mkdir -p /etc/mysql/conf.d/

WORKDIR /opt/mariadb

EXPOSE 3306

VOLUME /var/lib/mysql

# USER mysql

COPY scripts/docker-entrypoint.sh /opt/mariadb
ENTRYPOINT ["./docker-entrypoint.sh"]

# Default arguments passed to ENTRYPOINT if no arguments are passed when starting container
CMD ["mysqld_safe"]
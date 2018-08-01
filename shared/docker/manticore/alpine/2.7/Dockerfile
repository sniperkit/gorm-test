
FROM alpine:3.7 as builder

RUN apk add --no-cache \
    git \
    cmake \
    make \
    g++ \
    bison \
    flex \
    mariadb-dev \
    postgresql-dev \
    expat-dev \
    mariadb-client

RUN mkdir /build && cd /build \
    && git clone https://github.com/manticoresoftware/manticore.git  \
    && cd manticore && git checkout manticore-2.7.1 \
    && mkdir -p build && cd build \
    && cmake \
        -D SPLIT_SYMBOLS=1 \
        -D WITH_MYSQL=ON \
        -D WITH_PGSQL=ON \
        -D WITH_RE2=ON \
        -D WITH_STEMMER=ON \
        -D DISABLE_TESTING=ON \
        -D CMAKE_INSTALL_PREFIX=/usr \
        -D CONFFILEDIR=/etc/sphinxsearch \
        -D SPHINX_TAG=release .. \
    && make -j4 searchd indexer indextool

COPY config/sphinx.conf /build/manticore/build/src/

FROM alpine:3.7

RUN apk add --no-cache \
    mariadb-libs \
    mariadb-client-libs \
    postgresql-libs \
    expat \
    && mkdir -p /var/lib/manticore/log && mkdir -p /var/lib/manticore/data/

COPY --from=builder /build/manticore/build/src/indexer /usr/bin/
COPY --from=builder /build/manticore/build/src/indextool /usr/bin/
COPY --from=builder /build/manticore/build/src/searchd /usr/bin/
COPY --from=builder /build/manticore/build/src/sphinx.conf /etc/sphinxsearch/sphinx.conf

COPY scripts/entrypoint.sh /opt/manticore/entrypoint.sh
RUN chmod +x /opt/manticore/entrypoint.sh

VOLUME ["/var/lib/manticore", "/etc/sphinxsearch"]
EXPOSE 9306 9312 3306

ENTRYPOINT [ "/opt/manticore/entrypoint.sh" ]
CMD ["searchd", "--nodetach"]

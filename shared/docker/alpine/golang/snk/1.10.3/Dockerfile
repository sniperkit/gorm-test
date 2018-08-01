####################################
## Builder - image arguments
####################################
ARG BUILDER_ALPINE_VERSION=${BUILDER_ALPINE_VERSION:-"3.7"}
ARG BUILDER_GOLANG_VERSION=${BUILDER_GOLANG_VERSION:-"1.10.3"}
ARG BUILDER_IMAGE_TAG=${BUILDER_IMAGE_TAG:-"${BUILDER_GOLANG_VERSION}-alpine${BUILDER_ALPINE_VERSION}"}
ARG BUILDER_IMAGE_NAME=${BUILDER_IMAGE_NAME:-"golang:${BUILDER_IMAGE_TAG}"}

####################################
## Builder
###################################
FROM golang:1.10.3-alpine3.7

ARG APK_BUILDER=${APK_BUILDER:-"gcc g++ make ca-certificates openssl git cmake mercurial make nano bash jq musl-dev wget curl alpine-sdk sqlite-dev sqlite-libs sqlite tree"}

ENV PATH=${PATH:-"$PATH:$GOPATH/bin"}

RUN ALPINE_GLIBC_BASE_URL="https://github.com/sgerrand/alpine-pkg-glibc/releases/download" && \
    ALPINE_GLIBC_PACKAGE_VERSION="2.27-r0" && \
    ALPINE_GLIBC_BASE_PACKAGE_FILENAME="glibc-$ALPINE_GLIBC_PACKAGE_VERSION.apk" && \
    ALPINE_GLIBC_BIN_PACKAGE_FILENAME="glibc-bin-$ALPINE_GLIBC_PACKAGE_VERSION.apk" && \
    ALPINE_GLIBC_I18N_PACKAGE_FILENAME="glibc-i18n-$ALPINE_GLIBC_PACKAGE_VERSION.apk" && \
    apk add --no-cache --virtual=.build-dependencies wget ca-certificates && \
    wget \
        "https://raw.githubusercontent.com/sgerrand/alpine-pkg-glibc/master/sgerrand.rsa.pub" \
        -O "/etc/apk/keys/sgerrand.rsa.pub" && \
    wget \
        "$ALPINE_GLIBC_BASE_URL/$ALPINE_GLIBC_PACKAGE_VERSION/$ALPINE_GLIBC_BASE_PACKAGE_FILENAME" \
        "$ALPINE_GLIBC_BASE_URL/$ALPINE_GLIBC_PACKAGE_VERSION/$ALPINE_GLIBC_BIN_PACKAGE_FILENAME" \
        "$ALPINE_GLIBC_BASE_URL/$ALPINE_GLIBC_PACKAGE_VERSION/$ALPINE_GLIBC_I18N_PACKAGE_FILENAME" && \
    apk add --no-cache \
        "$ALPINE_GLIBC_BASE_PACKAGE_FILENAME" \
        "$ALPINE_GLIBC_BIN_PACKAGE_FILENAME" \
        "$ALPINE_GLIBC_I18N_PACKAGE_FILENAME" && \
    \
    rm "/etc/apk/keys/sgerrand.rsa.pub" && \
    /usr/glibc-compat/bin/localedef --force --inputfile POSIX --charmap UTF-8 "$LANG" || true && \
    echo "export LANG=$LANG" > /etc/profile.d/locale.sh && \
    \
    apk del glibc-i18n && \
    \
    rm "/root/.wget-hsts" && \
    apk del .build-dependencies && \
    rm \
        "$ALPINE_GLIBC_BASE_PACKAGE_FILENAME" \
        "$ALPINE_GLIBC_BIN_PACKAGE_FILENAME" \
        "$ALPINE_GLIBC_I18N_PACKAGE_FILENAME" \
    \
        && apk --no-cache add ${APK_BUILDER} \
    \
        && echo "GOPATH: $GOPATH" \
    \
        && go get -u github.com/vektah/gqlgen \
        && go install github.com/vektah/gqlgen \
    \
        && go get -u github.com/mitchellh/gox \
        && go install github.com/mitchellh/gox \
    \
        && go get -u github.com/Masterminds/glide \
        && go install github.com/Masterminds/glide \
    \
        && go get -u github.com/golang/dep/cmd/dep \
        && go install github.com/golang/dep/cmd/dep \
    \
        && go get -u github.com/mattn/gom \
        && go install github.com/mattn/gom \
    \
        && go get -u github.com/google/zoekt/... \
        && go install github.com/google/zoekt/cmd/... \
    \
        && go get -v -u github.com/kataras/bindata/cmd/... \
        && go install github.com/kataras/bindata/cmd/... \
    \
        && go get -u github.com/jteeuwen/go-bindata/... \
        && go install github.com/jteeuwen/go-bindata/... \
    \
        && go get -u github.com/svent/sift \
        && go install github.com/svent/sift \
    \
        && rm -fR $GOPATH/src \
        && rm -fR $GOPATH/pkg \
    \
        && ls -l $GOPATH/bin

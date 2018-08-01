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
FROM golang:${BUILDER_IMAGE_TAG} AS build

ARG REPO_VCS=${REPO_VCS:-"github.com"}
ARG REPO_NAMESPACE=${REPO_NAMESPACE:-"chrislusf"}
ARG REPO_PROJECT=${REPO_PROJECT:-"seaweedfs"}
ARG REPO_URI=${REPO_URI:-"${REPO_VCS}/${REPO_NAMESPACE}/${REPO_PROJECT}"}
ARG REPO_CLONE_DIR=${REPO_CLONE_DIR:-"/go/src/${REPO_URI}"}

ARG APK_BUILDER=${APK_BUILDER:-"gcc g++ make go musl-dev libgit2-dev wget curl ca-certificates openssl git cmake make nano mercurial bash"}

ENV LANG=${LANG:-"C.UTF-8"}

# Here we install GNU libc (aka glibc) and set C.UTF-8 locale as default.
RUN ALPINE_GLIBC_BASE_URL="https://github.com/sgerrand/alpine-pkg-glibc/releases/download" && \
    ALPINE_GLIBC_PACKAGE_VERSION="2.27-r0" && \
    ALPINE_GLIBC_BASE_PACKAGE_FILENAME="glibc-$ALPINE_GLIBC_PACKAGE_VERSION.apk" && \
    ALPINE_GLIBC_BIN_PACKAGE_FILENAME="glibc-bin-$ALPINE_GLIBC_PACKAGE_VERSION.apk" && \
    ALPINE_GLIBC_I18N_PACKAGE_FILENAME="glibc-i18n-$ALPINE_GLIBC_PACKAGE_VERSION.apk" && \
    apk add --no-cache --virtual=.build-dependencies ${APK_BUILDER} && \
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
    \
    go get -v -u ${REPO_URI}/weed && \
    go install ${REPO_URI}/weed && \
    \
    apk del .build-dependencies && \
    rm \
        "$ALPINE_GLIBC_BASE_PACKAGE_FILENAME" \
        "$ALPINE_GLIBC_BIN_PACKAGE_FILENAME" \
        "$ALPINE_GLIBC_I18N_PACKAGE_FILENAME"

############################################################################################################
############################################################################################################

####################################
## Builder - image arguments
####################################
ARG RUNNER_ALPINE_VERSION=${RUNNER_ALPINE_VERSION:-"3.7"}
ARG RUNNER_IMAGE_NAME=${RUNNER_IMAGE_NAME:-"alpine:${RUNNER_ALPINE_VERSION}"}

####################################
## Runner
####################################
FROM alpine:3.7

WORKDIR /usr/bin
COPY --from=build /go/bin .

EXPOSE 8080
EXPOSE 9333

VOLUME ["/data"]

COPY scripts/entrypoint.sh /opt/seaweedfs/entrypoint.sh
RUN chmod +x /opt/seaweedfs/entrypoint.sh

ENTRYPOINT ["/opt/seaweedfs/entrypoint.sh"]

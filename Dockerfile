####################################
## Builder - image arguments
####################################
# Note: issue with crane with build args as .env file variables...
#       crane: https://github.com/michaelsauter/crane
# ARG BUILDER_ALPINE_VERSION=${BUILDER_ALPINE_VERSION:-"3.7"}
# ARG BUILDER_GOLANG_VERSION=${BUILDER_GOLANG_VERSION:-"1.10.3"}
# ARG BUILDER_IMAGE_TAG=${BUILDER_IMAGE_TAG:-"${BUILDER_GOLANG_VERSION}-alpine${BUILDER_ALPINE_VERSION}"}
# ARG BUILDER_IMAGE_NAME=${BUILDER_IMAGE_NAME:-"golang:${BUILDER_IMAGE_TAG}"}

####################################
## Builder
###################################
FROM golang:1.10.3-alpine3.7 AS builder

ARG REPO_VCS=${REPO_VCS:-"github.com"}
ARG REPO_NAMESPACE=${REPO_NAMESPACE:-"sniperkit"}
ARG REPO_PROJECT=${REPO_PROJECT:-"gorm-test"}
ARG REPO_URI=${REPO_URI:-"${REPO_VCS}/${REPO_NAMESPACE}/${REPO_PROJECT}"}

## add apk build dependencies
RUN apk --no-cache --no-progress add gcc g++ make ca-certificates openssl git mercurial cmake make

WORKDIR /go/src/${REPO_URI}

## deps
COPY glide.lock glide.yaml ./
# COPY vendor vendor

## install commands
RUN go get -u github.com/Masterminds/glide \
    && glide install --strip-vendor

## pkg
COPY cmd cmd

## executables
# COPY cmd/meow cmd/meow
COPY cmd/gorm-load gorm-load

RUN go install ./... \
    && ls -la /go/bin

############################################################################################################
############################################################################################################

####################################
## Builder - image arguments
####################################
# Note: issue with crane with build args as .env file variables...
# ARG RUNNER_ALPINE_VERSION=${RUNNER_ALPINE_VERSION:-"3.7"}
# ARG RUNNER_IMAGE_NAME=${RUNNER_IMAGE_NAME:-"alpine:${RUNNER_ALPINE_VERSION}"}

####################################
## Build
####################################
FROM alpine:3.7 AS dist
WORKDIR /usr/bin

COPY --from=builder /go/bin ./

RUN echo "\n---- DEBUG INFO -----\n" \
    ls -l /usr/bin/gorm-* \
    echo "\nPATH: ${PATH}\n"

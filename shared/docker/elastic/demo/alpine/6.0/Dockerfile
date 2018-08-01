##########################################################
# DOWNLOADER                                             #
##########################################################
FROM alpine:3.7 as downloader

LABEL maintainer "https://github.com/blacktop"

ENV VERSION 6.0.1
ENV GOPATH /go

RUN apk add --no-cache ca-certificates \
                          build-base \
                          mercurial \
                          musl-dev \
                          openssl \
                          wget \
                          bash \
                          git \
                          go

RUN echo "===> Install filebeat..."
RUN mkdir -p ${GOPATH}/src/github.com/elastic
RUN cd ${GOPATH}/src/github.com/elastic && git clone -b v$VERSION --depth 1 https://github.com/elastic/beats.git
RUN cd ${GOPATH}/src/github.com/elastic/beats/filebeat \
  && go get \
  && make \
  && ./filebeat --version

RUN echo "===> Download nginx_logs..." \
  && wget --progress=bar:force -O /tmp/nginx_logs \
    https://raw.githubusercontent.com/elastic/examples/master/Common%20Data%20Formats/nginx_logs/nginx_logs
    # https://github.com/elastic/examples/raw/master/Common%20Data%20Formats/apache_logs/apache_logs

RUN ls -lah /tmp
RUN ls -lah ${GOPATH}/src/github.com/elastic/beats/filebeat
RUN ls -lah ${GOPATH}

##########################################################
# ES-DATA                                                #
##########################################################
FROM alpine:latest

# ENV BULK_SIZE 5242880

# RUN apk --no-cache add ca-certificates

# RUN mkdir -p /usr/share/filebeat/kibana/default
COPY --from=downloader /go/src/github.com/elastic/beats/filebeat/filebeat /bin/filebeat
COPY --from=downloader /go/src/github.com/elastic/beats/filebeat/filebeat.yml /usr/share/filebeat/filebeat.yml
COPY kibana/ /usr/share/filebeat/kibana/
COPY module/nginx /usr/share/filebeat/module/nginx
COPY modules.d/nginx.yml /usr/share/filebeat/modules.d/nginx.yml
COPY module/apache2 /usr/share/filebeat/module/apache2
COPY modules.d/apache2.yml /usr/share/filebeat/modules.d/apache2.yml
COPY fields.yml /usr/share/filebeat/fields.yml

COPY --from=downloader /tmp/nginx_logs /data/nginx_logs
COPY load.sh /load.sh

RUN ls -lah /usr/share/filebeat
RUN du -sh /usr/share/filebeat

WORKDIR /usr/share/filebeat

ENTRYPOINT ["/load.sh"]

FROM alpine:3.7
LABEL maintainer=l.stiebellehner@dodax.com

# run this container with `--volumes-from cassandra`

ENV MODE "snapshot"
ENV CASSANDRA_DATA_DIR="/var/lib/cassandra/data"
# leave keyspace empty to back up all
ENV CASSANDRA_KEYSPACE=""

ENV AWS_ACCESS_KEY_ID ""
ENV AWS_SECRET_ACCESS_KEY ""
ENV AWS_REGION "eu-central-1"
ENV S3_BUCKET ""

ENV TABLESNAP_NAME ""
# r = recursive, a = auto-add new subdirectories, B = Backup existing files to S3 if they are not already there
ENV TABLESNAP_OPTS "-r -a -B"
ENV TABLESNAP_THREADS "1"
ENV TABLESNAP_EXCLUDES "/snapshots/\|-tmp-\|cassandra.log"

# options for restore
# r = recursive
ENV TABLESLURP_OPTS "-r"

RUN apk add --no-cache python py-pip bash \
    && pip install awscli tablesnap

COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/bin/bash", "-c", "/entrypoint.sh"]
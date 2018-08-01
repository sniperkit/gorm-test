FROM srakshit/alpine-base

MAINTAINER Subham Rakshit

ENV CASSANDRA_HOME="/opt/cassandra" \
    CASSANDRA_VERSION="2.1.17" \
    PATH=$PATH:$CASSANDRA_HOME/bin

ADD ./cassandra-* /opt/

RUN apk add --no-cache python python-dev

RUN chmod +x /opt/cassandra-download.sh \
    && sleep 1 \
    && /opt/cassandra-download.sh \
    && rm /opt/cassandra-download.sh \
    && mv /opt/apache-cassandra-$CASSANDRA_VERSION /opt/cassandra-$CASSANDRA_VERSION \
    && ln -s /opt/cassandra-$CASSANDRA_VERSION $CASSANDRA_HOME \
    && addgroup cassandra \
    && adduser -S -G cassandra cassandra \
    && chown cassandra:cassandra -R /opt/cassandra* \
    && chmod 754 /opt/cassandra*

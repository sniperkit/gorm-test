FROM alpine:3.7

ARG CASSANDRA_VERSION="3.11.2"

# set cassandra directories
ENV CASSANDRA_HOME=/opt/cassandra
ENV CASSANDRA_CONFIG=/opt/cassandra/conf
ENV CASSANDRA_PERSIST_DIR=/var/lib/cassandra
ENV CASSANDRA_DATA=/var/lib/cassandra/data
ENV CASSANDRA_LIB=/usr/share/cassandra/lib
ENV CASSANDRA_COMMITLOG=/var/lib/cassandra/commitlog
ENV CASSANDRA_LOG=/var/lib/cassandra/log
ENV CASSANDRA_USER=cassandra
ENV CASSANDRA_PEER_DISCOVERY_SERVICE="cassandra-peers"

ENV JVM_OPTS "-XX:+UseGCLogFileRotation -Xloggc:/tmp/cassandra-gc.log -XX:NumberOfGCLogFiles=1 -XX:GCLogFileSize=10M -javaagent:/usr/share/cassandra/lib/jolokia-agent.jar=port=8778,host=0.0.0.0"
ENV HEAP_SETTINGS "-Xms4096M -Xmx4096M"

RUN apk add --no-cache openjdk8-jre jemalloc procps iproute2 gnupg ca-certificates wget curl \
    && wget -O /usr/local/bin/gosu "https://github.com/tianon/gosu/releases/download/1.10/gosu-amd64" \
    && chmod +x /usr/local/bin/gosu \
    && gosu nobody true \
    && addgroup -g 998 cassandra \
    && adduser -D -H -u 998 -G cassandra cassandra

## Create data directories that should be used by Cassandra
RUN mkdir -p ${CASSANDRA_HOME} \
    ${CASSANDRA_DATA} \
    ${CASSANDRA_CONFIG} \
    ${CASSANDRA_LOG} \
    ${CASSANDRA_COMMITLOG} \
    ${CASSANDRA_LIB}

# install Cassandra
RUN apk --update --no-cache add wget ca-certificates tar bash shadow bind-tools && \
    wget https://archive.apache.org/dist/cassandra/${CASSANDRA_VERSION}/apache-cassandra-${CASSANDRA_VERSION}-bin.tar.gz -P /tmp && \
    tar -xvzf /tmp/apache-cassandra-${CASSANDRA_VERSION}-bin.tar.gz -C ${CASSANDRA_HOME} --strip-components=1 && \
    apk --purge del wget ca-certificates tar && \
    rm -r /tmp/apache-cassandra-${CASSANDRA_VERSION}-bin.tar.gz \
    /var/cache/apk/*

# install jolokia agent
RUN wget -q https://repo1.maven.org/maven2/org/jolokia/jolokia-jvm/1.3.7/jolokia-jvm-1.3.7-agent.jar -O ${CASSANDRA_LIB}/jolokia-agent.jar \
    && chmod +x ${CASSANDRA_LIB}/jolokia-agent.jar

# Add default config
RUN mv ${CASSANDRA_HOME}/conf/* ${CASSANDRA_CONFIG}
COPY ./conf/* ${CASSANDRA_CONFIG}/

# https://issues.apache.org/jira/browse/CASSANDRA-11661
RUN chmod +x ${CASSANDRA_CONFIG}/*.sh \
    && sed -ri 's/^(JVM_PATCH_VERSION)=.*/\1=25/' ${CASSANDRA_CONFIG}/cassandra-env.sh

# Add cassandra bin to PATH
ENV PATH=$PATH:${CASSANDRA_HOME}/bin \
    CASSANDRA_CONF=${CASSANDRA_CONFIG}

# Change directory ownership and access rights
RUN chown -R ${CASSANDRA_USER}:${CASSANDRA_USER} \
    ${CASSANDRA_HOME} \
    ${CASSANDRA_PERSIST_DIR} \
    ${CASSANDRA_LOG} \
    ${CASSANDRA_CONFIG} \
    ${CASSANDRA_LIB}

# Expose data volume
VOLUME ${CASSANDRA_PERSIST_DIR}

COPY entrypoint.sh /entrypoint.sh
COPY repair-cassandra.sh /repair-cassandra.sh
RUN chmod +x /entrypoint.sh /repair-cassandra.sh

#USER ${CASSANDRA_USER}
WORKDIR ${CASSANDRA_HOME}

# 7000: intra-node communication
# 7001: TLS intra-node communication
# 7199: JMX
# 8788: jolokia metrics
# 9042: CQL
EXPOSE 7000 7001 7199 8788 9042

ENTRYPOINT ["/bin/bash", "/entrypoint.sh"]

CMD ["cassandra", "-f"]
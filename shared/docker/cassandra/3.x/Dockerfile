FROM cassandra:3.11
LABEL maintainer="l.stiebellehner@dodax.com"

ENV JVM_OPTS "-Xms4096M -Xmx4096M -javaagent:/usr/share/cassandra/lib/jolokia-agent.jar=port=8778,host=0.0.0.0"

VOLUME /var/lib/cassandra

COPY repair-cassandra.sh /opt/repair-cassandra.sh
COPY snapshot-cassandra.sh /opt/snapshot/cassandra.sh

RUN apt-get -qq update \
    && apt-get -qq install wget python3-pip \
    && wget -q https://repo1.maven.org/maven2/org/jolokia/jolokia-jvm/1.3.7/jolokia-jvm-1.3.7-agent.jar -O /usr/share/cassandra/lib/jolokia-agent.jar \
    && chmod +x /usr/share/cassandra/lib/jolokia-agent.jar /opt/repair-cassandra.sh /opt/snapshot-cassandra.sh \
    && apt-get -qq -y purge wget \
    && apt-get -qq -y autoremove \
    && pip install --upgrade pip \
    && pip install awscli

ENTRYPOINT ["/docker-entrypoint.sh"]

# 7000: intra-node communication
# 7001: TLS intra-node communication
# 7199: JMX
# 8778: Jolokia
# 9042: CQL
# 9160: thrift service
EXPOSE 7000 7001 7199 8778 9042 9160

CMD ["cassandra", "-f"]

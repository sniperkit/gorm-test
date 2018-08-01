ARG HADOOP_VERSION=2.9.0

FROM mpolatcan/hadoop:ubuntu-${HADOOP_VERSION}

ARG HIVE_VERSION=2.3.3
ARG HIVE_DOWNLOAD_LINK="http://www-eu.apache.org/dist/hive/hive-${HIVE_VERSION}/apache-hive-${HIVE_VERSION}-bin.tar.gz"

# TODO Configurations will be added

MAINTAINER Mutlu Polatcan <mutlupolatcan@gmail.com>

ENV HIVE_HOME "/usr/local/hive"
ENV PATH $PATH:$HIVE_HOME/bin

COPY entrypoint.sh hive_config_loader.sh ./

RUN wget ${HIVE_DOWNLOAD_LINK} && \
    mkdir -p ${HIVE_HOME} && \
    tar -xvzf apache-hive-${HIVE_VERSION}-bin.tar.gz -C /usr/local/ && \
    mv /usr/local/apache-hive-${HIVE_VERSION}-bin/* ${HIVE_HOME} && \
    rm -r /usr/local/apache-hive-${HIVE_VERSION}-bin && \
    rm apache-hive-${HIVE_VERSION}-bin.tar.gz && \
    chown -R hduser:hadoop ${HIVE_HOME} && \
    chmod -R 777 ${HIVE_HOME} && \
    echo "export HADOOP_HOME=/usr/local/hadoop" >> ${HIVE_HOME}/conf/hive-env.sh && \
    chmod +x entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]
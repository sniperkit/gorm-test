ARG NEO4J_BASE_IMAGE_NAME=${NEO4J_BASE_IMAGE_NAME:-"openjdk"}
ARG NEO4J_BASE_IMAGE_TAG=${NEO4J_BASE_IMAGE_TAG:-"8-jre-alpine"}
ARG NEO4J_BASE_IMAGE=${NEO4J_BASE_IMAGE:-"${NEO4J_BASE_IMAGE_NAME}:${NEO4J_BASE_IMAGE_TAG}"}

# fix crane env file error
# FROM ${NEO4J_BASE_IMAGE}
FROM openjdk:8-jre-alpine

# neo4j - version
ARG NEO4J_EDITION=${NEO4J_EDITION:-"community"}
ARG NEO4J_VERSION=${NEO4J_VERSION:-"3.4.5"}
ARG NEO4J_OS=${NEO4J_OS:-"unix"}

ARG NEO4J_SHA256=${NEO4J_SHA256:-"af53823776645e11d04436a513368e7e417b515572d6228da6b2977c8490ffbb"}
ARG NEO4J_TARBALL=${NEO4J_TARBALL:-"neo4j-${NEO4J_EDITION}-${NEO4J_VERSION}-${NEO4J_OS}.tar.gz"}
ARG NEO4J_URI=${NEO4J_URI:-"http://dist.neo4j.org/neo4j-${NEO4J_EDITION}-${NEO4J_VERSION}-${NEO4J_OS}.tar.gz"}

# neo4j - user/group
ARG NEO4J_USER=${NEO4J_USER:-"neo4j"}
ARG NEO4J_GROUP=${NEO4J_GROUP:-"neo4j"}
ARG NEO4J_HOME=${NEO4J_HOME:-"/var/lib/neo4j"}

# neo4j - output directories
ARG NEO4J_CERTS_DIR=${NEO4J_CERTS_DIR:-"/opt/neo4j/ssl"}
ARG NEO4J_CONF_DIR=${NEO4J_CONF_DIR:-"/opt/neo4j/conf"}
ARG NEO4J_DATA_DIR=${NEO4J_DATA_DIR:-"/opt/neo4j/data"}
ARG NEO4J_LOGS_DIR=${NEO4J_LOGS_DIR:-"/opt/neo4j/logs"}
ARG NEO4J_IMPORT_DIR=${NEO4J_IMPORT_DIR:-"/opt/neo4j/import"}
ARG NEO4J_PLUGINS_DIR=${NEO4J_PLUGINS_DIR:-"/opt/neo4j/plugins"}
ARG NEO4J_METRICS_DIR=${NEO4J_METRICS_DIR:-"/opt/neo4j/metrics"}

# neo4j - command(s)
ARG NEO4J_CMD_DEFAULT=${NEO4J_CMD_DEFAULT:-"neo4j"}
ARG NEO4J_CMD_CONSOLE=${NEO4J_CMD_CONSOLE:-"console"}

# neo4j - apk dependencies
ARG NEO4J_APK_BUILD=${NEO4J_APK_BUILD:-"openssl ca-certificates curl"}
ARG NEO4J_APK_RUNTIME=${NEO4J_APK_RUNTIME:-"bash tini su-exec"}

COPY ./local-package/* /tmp/

RUN addgroup -S ${NEO4J_GROUP} \
    && adduser -S -H -h ${NEO4J_HOME} -G ${NEO4J_GROUP} ${NEO4J_USER} \
    \
    && apk add --no-cache --no-progress --quiet ${NEO4J_APK_BUILD} ${NEO4J_APK_RUNTIME} \
    \
    && curl --fail --silent --show-error --location --remote-name ${NEO4J_URI} \
    && echo "${NEO4J_SHA256}  ${NEO4J_TARBALL}" | sha256sum -csw - \
    && tar --extract --file ${NEO4J_TARBALL} --directory /var/lib \
    && mv ${NEO4J_HOME}-* ${NEO4J_HOME} \
    && rm ${NEO4J_TARBALL} \
    \
    && mv ${NEO4J_HOME}/data ${NEO4J_DATA_DIR} \
    && chown -R ${NEO4J_GROUP}:${NEO4J_USER} ${NEO4J_DATA_DIR} \
    && chmod -R 777 ${NEO4J_DATA_DIR} \
    && chown -R ${NEO4J_GROUP}:${NEO4J_USER} ${NEO4J_HOME} \
    && chmod -R 777 ${NEO4J_HOME} \
    && ln -s ${NEO4J_DATA_DIR} ${NEO4J_HOME}/data \
    \
    && apk del ${NEO4J_APK_BUILD}

ENV PATH ${PATH:-"${NEO4J_HOME}/bin:$PATH"} \
    EXTENSION_SCRIPT=${EXTENSION_SCRIPT:-""} \
    NEO4J_ACCEPT_LICENSE_AGREEMENT=${NEO4J_ACCEPT_LICENSE_AGREEMENT:-"yes"} \
    NEO4J_AUTH=${NEO4J_AUTH:-"default:secret"}

WORKDIR ${NEO4J_HOME}

VOLUME ["${NEO4J_DATA_DIR}"]

COPY ./scripts/docker-entrypoint.sh /opt/scripts/docker-entrypoint.sh

EXPOSE 7474 7473 7687

ENTRYPOINT ["/sbin/tini", "-g", "--", "/opt/scripts/docker-entrypoint.sh"]
CMD ["${NEO4J_CMD_DEFAULT}"]

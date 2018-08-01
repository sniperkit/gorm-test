#!/bin/bash
set -e

export ELASTICSEARCH_URL=${ELASTICSEARCH_URL:-"http://localhost:9200"}

echo ""
echo " - ARG_0: $0"
echo " - ARG_1: $1"
echo " - ARG_2: $2"
echo " - ARGS: $@"
echo " - ES_JAVA_OPTS: ${ES_JAVA_OPTS}"

echo ""
echo " - ELSK_USER: ${ELSK_USER}"
echo " - ELSK_PASS: ${ELSK_PASS}"
echo " - ELSK_DOMAIN: ${ELSK_DOMAIN}"
echo " - ELASTICSEARCH_URL: ${ELASTICSEARCH_URL}"
echo ""

# Add kibana as command if needed
if [[ "$1" == -* ]]; then
	set -- kibana "$@"
fi

# Run as user "elstack" if the command is "kibana"
if [ "$1" = 'kibana' ]; then
	if [ "$ELASTICSEARCH_URL" ]; then
		sed -ri "s!^(\#\s*)?(elasticsearch\.url:).*!\2 '$ELASTICSEARCH_URL'!" /usr/share/kibana/config/kibana.yml
        # /etc/kibana/kibana.yml
	fi
    echo " - cat /etc/kibana/kibana.yml: "
    cat /etc/kibana/kibana.yml
	set -- su-exec elstack tini -- "$@"
fi

exec "$@"

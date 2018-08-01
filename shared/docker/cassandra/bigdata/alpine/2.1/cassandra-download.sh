#!/bin/bash

echo "Downloading and unzipping apache cassandra-$CASSANDRA_VERSION"
mirror=$(curl --stderr /dev/null https://www.apache.org/dyn/closer.cgi\?as_json\=1 | jq -r '.preferred')
url="$mirror/cassandra/$CASSANDRA_VERSION/apache-cassandra-${CASSANDRA_VERSION}-bin.tar.gz"
curl $url | tar xz -C /opt/

#!/bin/sh

set -e

# filebeat setup -e -E output.elasticsearch.hosts=elasticsearch:9200 -E setup.kibana.host=kibana:5601

# Ingest data into Elasticsearch
filebeat -e --modules=nginx --setup \
         -E output.elasticsearch.hosts=elasticsearch:9200 \
         -E setup.kibana.host=elasticsearch:80

# filebeat setup -e -E output.elasticsearch.hosts=localhost:9200 -E setup.kibana.host=localhost:5601
#
# filebeat -e --modules=nginx \
#          -E output.elasticsearch.hosts=localhost:9200 \
#          -E setup.kibana.host=localhost:5601

#!/bin/bash

if [ "$MODE" == "snapshot" ] ; then
    tablesnap ${TABLESNAP_OPTS} -t "${TABLESNAP_THREADS}" -n "${TABLESNAP_NAME}" \
    --exclude="${TABLESNAP_EXCLUDES}" --exclude="/system.*" --aws-region "${AWS_REGION}" \
    -k "${AWS_ACCESS_KEY_ID}" -s "${AWS_SECRET_ACCESS_KEY}" "${S3_BUCKET}" "${CASSANDRA_DATA_DIR}/${CASSANDRA_KEYSPACE}"
elif [ "$MODE" == "restore" ] ; then
    tableslurp ${TABLESLURP_OPTS} -t "${TABLESNAP_THREADS}" -n "${TABLESNAP_NAME}" --owner cassandra --group cassandra \
    --aws-region "${AWS_REGION}" -k "${AWS_ACCESS_KEY_ID}" -s "${AWS_SECRET_ACCESS_KEY}" -t "${THREADS}" "${S3_BUCKET}" \
    "${CASSANDRA_DATA_DIR}/${CASSANDRA_KEYSPACE}" "${CASSANDRA_DATA_DIR}/${CASSANDRA_KEYSPACE}"
else
    printf "%s\n" "Unknown MODE specified: ${MODE}, exiting."
    exit 1
fi

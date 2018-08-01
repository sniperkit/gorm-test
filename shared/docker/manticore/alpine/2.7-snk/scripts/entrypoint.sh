#!/bin/sh

set -x
set -e

if [ "${1#-}" != "$1" ]; then
    set -- searchd "$@"
fi

if [ "$1" = 'searchd' ]; then
  # If the sphinx config exists, try to run the indexer before starting searchd
  if [ -f /etc/sphinx/sphinx.conf ]; then
    indexer --all
  fi
fi

echo "check mariadb processes"
ps aux | grep mariadb

echo "check mysql processes"
ps aux | grep mysql

exec "$@"
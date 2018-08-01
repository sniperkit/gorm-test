#!/bin/sh

set -x
set -e

# WIP, not working...

case "$1" in

  'indexer')
    ARGS="-ip `hostname -i` -mdir /data"
  # Is this instance linked with an other master? (Docker commandline "--link master1:master")
  if [ -n "$MASTER_PORT_9333_TCP_ADDR" ] ; then
    ARGS="$ARGS -peers=$MASTER_PORT_9333_TCP_ADDR:$MASTER_PORT_9333_TCP_PORT"
  fi
    exec /usr/bin/weed $@ $ARGS
  ;;

  'searcher')
    ARGS="-ip `hostname -i` -dir /data"
  # Is this instance linked with a master? (Docker commandline "--link master1:master")
    if [ -n "$MASTER_PORT_9333_TCP_ADDR" ] ; then
    ARGS="$ARGS -mserver=$MASTER_PORT_9333_TCP_ADDR:$MASTER_PORT_9333_TCP_PORT"
  fi
    exec /usr/bin/weed $@ $ARGS
  ;;

  'daemon')
    ARGS="-ip `hostname -i` -dir /data"
    if [ -n "$MASTER_PORT_9333_TCP_ADDR" ] ; then
    ARGS="$ARGS -master.peers=$MASTER_PORT_9333_TCP_ADDR:$MASTER_PORT_9333_TCP_PORT"
  fi
    exec /usr/bin/weed $@ $ARGS
    ;;

  *)
    exec /usr/bin/weed $@
  ;;
esac
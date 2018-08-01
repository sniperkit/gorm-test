#!/bin/sh

# set -e
# set -x

# unalias stop
clear

NGINX_CMD=$(which nginx)
NGINX_CONF="/etc/nginx/snk.conf"
RETVAL=0

check_dirs() {
  ls -la /etc/nginx
  tree /etc/nginx
  tree /var/www/*
  tree /var/www/html/meow-vuejs-front
  cat /etc/hosts
  cat /etc/nginx/snk.conf
  cat /etc/hostname
}

check_endpoints() {
  echo " [PING] Nginx reverse-proxy"
  ping -c 2 nginx:8080
  echo ""

  echo " [PING] Kibana"
  ping -c 2 kibana:5601
  echo ""

  echo " [PING] Elstack"
  ping -c 2 elstack:8089
  echo ""

  echo " [PING] Front server"
  ping -c 2 front:3000
  echo ""

  echo " [PING] Meow QUERY service"
  ping -c 2 query:8080
  echo ""

  echo " [PING] Meow PUSHER service"
  ping -c 2 pusher:8080
  echo ""

  echo " [PING] Meow query service"
  ping -c 2 meow:8080
  echo ""

}

check_elastic_node() {
  echo " [PING] Master container"
  ping -c 2 master:9200
  echo ""
}

daemon() {
   echo "Starting NGINX Web Server as daemon:"
   $NGINX_CMD -c $NGINX_CONF &
   RETVAL=$?
   [ $RETVAL -eq 0 ] && echo "status: [OK]" || echo "status: [FAILED]"
   return $RETVAL
}

start() {
   echo "Starting NGINX Web Server as command..."
   $NGINX_CMD -c $NGINX_CONF &
   RETVAL=$?
   [ $RETVAL -eq 0 ] && echo "status: [OK]" || echo "status: [FAILED]"
   return $RETVAL
}

stop() {
   echo "Stopping NGINX Web Server..."
   $NGINX_CMD -s quit
   RETVAL=$?
   [ $RETVAL -eq 0 ] && echo "status: [OK]" || echo "status: [FAILED]"
   return $RETVAL
}

watch() {
  while true
  do
    clear
    echo " 1. Test/Validate NGINX Config... (luc3)"
    # Check NGINX Configuration Test
    # inotifywait --exclude .swp -e create -e modify -e delete -e move /etc/nginx/sites-enabled /etc/nginx
    inotifywait --exclude .swp -e create -e modify -e delete -e move /etc/nginx/snk.conf
    # Only Reload NGINX If NGINX Configuration Test Pass
    nginx -t
    if [ $? -eq 0 ]
    then
      clear
      # echo " 1a. Check directories --- "
      # check_dirs
      check_endpoints
      echo " 2. Reloading Nginx Configuration --- "
      nginx -s reload
    fi
  done
  return $RETVAL
}

echo ""
echo " Local config:"
echo " - NGINX_CMD: ${NGINX_CMD}"
echo " - NGINX_CONF: ${NGINX_CONF}"
echo ""

case "$1" in
  'start')
      start
    ;;

  'stop')
      stop
    ;;

  'restart')
      stop
      start
    ;;

  'watch')
      echo "Watching for NGINX configuration's changes..."
      start
      check_dirs
      watch
    ;;

  'bash')
      exec /bin/bash
    ;;

  'proxy')
    ARGS=""
    # ARGS="-ip `hostname -i` -dir /data"
    # if [ -n "$MASTER_PORT_9333_TCP_ADDR" ] ; then
    #   ARGS="$ARGS -master.peers=$MASTER_PORT_9333_TCP_ADDR:$MASTER_PORT_9333_TCP_PORT"
    # fi
    # exec /usr/bin/weed $@ $ARGS
    exec $@ $ARGS
    ;;

  'help')
    echo "Usage: $0 {start|stop|restart|watch}"
    exit 1
    ;;

   *)

  # default: start nginx with daemon mode disabled
  # exec nginx -g "daemon off;"
  exec nginx -g "daemon off;"
  echo "goodbye nginx..."
  exit 1

esac

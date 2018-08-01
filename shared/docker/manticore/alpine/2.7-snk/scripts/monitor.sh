#!/bin/bash
LOCKFILE="/tmp/sphinx_monitor.lock"
lockfile -r 0 $LOCKFILE || exit 1
RESTART=0
MIN=20000
DATE=`date`
QUERY="select count(*) as c from listing2 where match('100')"
IP=`/sbin/ifconfig | sed -En 's/127.0.0.1//;s/.*inet (addr:)?(([0-9]*\.){3}[0-9]*).*/\2/p'`
PORT="9309"

ITEMS=`timeout 5s mysql -B -h $IP -P $PORT -u test -e "$QUERY" -A -r -N`

if [ $? == '124' ]
then

 echo $DATE 'Timeout'
 RESTART=1
fi

if [ "$ITEMS" == "" ]
then

 echo $DATE 'Empty result'
 RESTART=1

else

 if [ "$ITEMS" -lt "$MIN" ]
 then

  echo $DATE Less then $MIN items
  RESTART=1

 fi
fi

if [ "$RESTART" -eq "1" ]
then

 echo $DATE Restarting sphinx
 /etc/init.d/sphinxsearch stop
 sleep 10
 if pgrep -x "searchd"
 then
  echo killing all searchd
  killall -s SIGKILL searchd
  sleep 2
 fi
 /etc/init.d/sphinxsearch start
 sleep 180

else

 echo $DATE $ITEMS items in sphinx $IP >> /var/log/sphinx_monitor.log

fi

rm -f $LOCKFILE

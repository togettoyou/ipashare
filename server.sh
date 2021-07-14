#!/bin/bash

# 启动服务名
SERVER_NAME="super-signature-app"
BASE_DIR=$PWD
INTERVAL=2

function start() {
  if [ "$(pgrep $SERVER_NAME -u $UID)" != "" ]; then
    echo "$SERVER_NAME already running"
    exit 1
  fi

  nohup "$BASE_DIR"/$SERVER_NAME &>/dev/null &

  echo "sleeping..." && sleep $INTERVAL

  if [ "$(pgrep $SERVER_NAME -u $UID)" == "" ]; then
    echo "$SERVER_NAME start failed"
    exit 1
  fi
}

function status() {
  if [ "$(pgrep $SERVER_NAME -u $UID)" != "" ]; then
    echo $SERVER_NAME is running
  else
    echo $SERVER_NAME is not running
  fi
}

function stop() {
  if [ "$(pgrep $SERVER_NAME -u $UID)" != "" ]; then
    kill -9 "$(pgrep $SERVER_NAME -u $UID)"
  fi

  echo "sleeping..." && sleep $INTERVAL

  if [ "$(pgrep $SERVER_NAME -u $UID)" != "" ]; then
    echo "$SERVER_NAME stop failed"
    exit 1
  fi
}

case "$1" in
'start')
  start
  ;;
'stop')
  stop
  ;;
'status')
  status
  ;;
'restart')
  stop && start
  ;;
*)
  echo "usage: $0 {start|stop|restart|status}"
  exit 1
  ;;
esac

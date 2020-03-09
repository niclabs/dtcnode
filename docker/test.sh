#!/usr/bin/env bash
set -e
cd "$(dirname $0)" || exit
go run github.com/niclabs/dtcconfig rsa \
  -n 0.0.0.0:9871,0.0.0.0:9873,0.0.0.0:9875,0.0.0.0:9877,0.0.0.0:9879 \
  -t 3 \
  -H "$(ip addr | grep 'global docker0' | awk '{print $2}' | sed sx/16xxg)" \
  -c "dtc-config.yaml" \
  -k "config/" \
  -d "/tmp/dtc.sqlite3"
docker-compose build
docker-compose up

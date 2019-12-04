#!/usr/bin/env bash
set -e
go build github.com/niclabs/dtcconfig rsa \
            -n 0.0.0.0:2030,0.0.0.0:2030,0.0.0.0:2030,0.0.0.0:2030,0.0.0.0:2030 \
            -t 3 \
            -H "$(ip addr | grep 'global docker0' | awk '{print $2}' | sed sx/16xxg)" \
            -c "config.yaml" \
            -k "config/"
docker-compose build
docker-compose up

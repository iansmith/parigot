#!/usr/bin/bash
set -e 

chmod 777 /var/run/docker-host.sock
docker info
echo "docker from docker working ok"
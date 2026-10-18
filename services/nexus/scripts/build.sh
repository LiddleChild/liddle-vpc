#!/usr/bin/env bash

set -e

docker build -t bosstnp/nexus:latest -f ../Dockerfile ..

docker push bosstnp/nexus:latest

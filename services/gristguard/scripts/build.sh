#!/usr/bin/env bash

set -e

docker build -t bosstnp/gristguard:latest -f ../Dockerfile ..

docker push bosstnp/gristguard:latest

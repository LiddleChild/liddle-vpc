#!/usr/bin/env bash

set -e

docker build -t bosstnp/backend:latest -f ../Dockerfile ..

docker push bosstnp/backend:latest

#!/usr/bin/env bash

set -e

docker build -t bosstnp/expenses:latest -f ../Dockerfile ..

docker push bosstnp/expenses:latest

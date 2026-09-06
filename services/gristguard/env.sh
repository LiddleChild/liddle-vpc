#!/usr/bin/env sh

export $(grep -v '^#' .env | xargs -0)


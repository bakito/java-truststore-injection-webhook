#!/bin/bash
set -e
docker build -f Dockerfile -t jti:e2e .
kind load docker-image jti:e2e --name e2e

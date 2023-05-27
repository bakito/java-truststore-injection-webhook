#!/bin/bash
set -e
docker build -f Dockerfile -t localhost:5001/jti:e2e .
docker push localhost:5001/jti:e2e

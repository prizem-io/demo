#!/bin/sh
docker build -t prizem-io/demo-frontend -f Dockerfile.frontend .
docker build -t prizem-io/demo-backend -f Dockerfile.backend .
docker build -t prizem-io/demo-message -f Dockerfile.message .
docker build -t prizem-io/demo-message-grpc -f Dockerfile.message-grpc .
docker build -t prizem-io/demo-register -f Dockerfile.register .
docker build -t prizem-io/demo-announce -f Dockerfile.announce .
docker build -t prizem-io/demo-proxy-init -f Dockerfile.proxy-init .